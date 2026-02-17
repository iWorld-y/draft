package biz

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"backend/internal/biz/entity"
	"backend/internal/biz/repo"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 7 * 24 * time.Hour

var (
	ErrInvalidCredentials = kerrors.Unauthorized("INVALID_CREDENTIALS", "用户名或密码错误")
	ErrUnauthorized       = kerrors.Unauthorized("UNAUTHORIZED", "未授权")
	ErrUserExists         = kerrors.BadRequest("USER_EXISTS", "用户已存在")
	ErrInvalidInput       = kerrors.BadRequest("INVALID_INPUT", "用户名和密码不能为空")
)

type AuthUseCase struct {
	userRepo    repo.UserRepo
	tokenRepo   repo.RefreshTokenRepo
	jwtSecret   []byte
	accessToken time.Duration
	refreshTTL  time.Duration
}

func NewAuthUseCase(userRepo repo.UserRepo, tokenRepo repo.RefreshTokenRepo) *AuthUseCase {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		secret = "draft-dev-secret-change-me"
	}

	return &AuthUseCase{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		jwtSecret:   []byte(secret),
		accessToken: tokenTTL,
		refreshTTL:  tokenTTL,
	}
}

type AuthResult struct {
	User         *entity.User
	AccessToken  string
	RefreshToken string
}

func (uc *AuthUseCase) Register(ctx context.Context, username, password string) (*AuthResult, error) {
	normalizedUsername := strings.TrimSpace(username)
	if normalizedUsername == "" || password == "" {
		return nil, ErrInvalidInput
	}

	existing, err := uc.userRepo.GetByUsername(ctx, normalizedUsername)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		Username:     normalizedUsername,
		PasswordHash: string(hashed),
		Status:       1,
	}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return uc.issueSession(ctx, user)
}

func (uc *AuthUseCase) Login(ctx context.Context, username, password string) (*AuthResult, error) {
	normalizedUsername := strings.TrimSpace(username)
	if normalizedUsername == "" || password == "" {
		return nil, ErrInvalidInput
	}

	user, err := uc.userRepo.GetByUsername(ctx, normalizedUsername)
	if err != nil {
		return nil, err
	}
	if user == nil || user.DeletedAt != nil || user.Status != 1 {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return uc.issueSession(ctx, user)
}

func (uc *AuthUseCase) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, ErrUnauthorized
	}

	tokenHash := hashToken(refreshToken)
	existing, err := uc.tokenRepo.GetValidByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrUnauthorized
	}

	user, err := uc.userRepo.GetByID(ctx, existing.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.DeletedAt != nil || user.Status != 1 {
		return nil, ErrUnauthorized
	}

	if err := uc.tokenRepo.RevokeByHash(ctx, tokenHash); err != nil {
		return nil, err
	}

	return uc.issueSession(ctx, user)
}

func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return nil
	}
	return uc.tokenRepo.RevokeByHash(ctx, hashToken(refreshToken))
}

func (uc *AuthUseCase) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.DeletedAt != nil || user.Status != 1 {
		return nil, ErrUnauthorized
	}
	return user, nil
}

func (uc *AuthUseCase) ParseAccessToken(accessToken string) (int64, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnauthorized
		}
		return uc.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrUnauthorized
	}

	subVal, ok := claims["sub"]
	if !ok {
		return 0, ErrUnauthorized
	}

	subFloat, ok := subVal.(float64)
	if !ok {
		return 0, ErrUnauthorized
	}

	userID := int64(subFloat)
	if userID <= 0 {
		return 0, ErrUnauthorized
	}
	return userID, nil
}

func (uc *AuthUseCase) issueSession(ctx context.Context, user *entity.User) (*AuthResult, error) {
	accessToken, err := uc.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRandomToken(32)
	if err != nil {
		return nil, err
	}

	rt := &entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: time.Now().Add(uc.refreshTTL),
	}
	if err := uc.tokenRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *AuthUseCase) generateAccessToken(user *entity.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"iat":      now.Unix(),
		"exp":      now.Add(uc.accessToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return signed, nil
}

func generateRandomToken(byteLen int) (string, error) {
	buf := make([]byte, byteLen)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
