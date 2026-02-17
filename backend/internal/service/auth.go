package service

import (
	"context"

	"backend/internal/biz"
)

const RefreshCookieName = "refresh_token"

type AuthService struct {
	uc *biz.AuthUseCase
}

func NewAuthService(uc *biz.AuthUseCase) *AuthService {
	return &AuthService{uc: uc}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type AuthResponse struct {
	AccessToken string   `json:"access_token"`
	User        AuthUser `json:"user"`
}

func (s *AuthService) Register(ctx context.Context, req *AuthRequest) (*AuthResponse, string, error) {
	result, err := s.uc.Register(ctx, req.Username, req.Password)
	if err != nil {
		return nil, "", err
	}

	return &AuthResponse{
		AccessToken: result.AccessToken,
		User: AuthUser{
			ID:       result.User.ID,
			Username: result.User.Username,
		},
	}, result.RefreshToken, nil
}

func (s *AuthService) Login(ctx context.Context, req *AuthRequest) (*AuthResponse, string, error) {
	result, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, "", err
	}

	return &AuthResponse{
		AccessToken: result.AccessToken,
		User: AuthUser{
			ID:       result.User.ID,
			Username: result.User.Username,
		},
	}, result.RefreshToken, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, string, error) {
	result, err := s.uc.Refresh(ctx, refreshToken)
	if err != nil {
		return nil, "", err
	}

	return &AuthResponse{
		AccessToken: result.AccessToken,
		User: AuthUser{
			ID:       result.User.ID,
			Username: result.User.Username,
		},
	}, result.RefreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.uc.Logout(ctx, refreshToken)
}

func (s *AuthService) Me(ctx context.Context, userID int64) (*AuthUser, error) {
	user, err := s.uc.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &AuthUser{ID: user.ID, Username: user.Username}, nil
}

func (s *AuthService) ParseAccessToken(accessToken string) (int64, error) {
	return s.uc.ParseAccessToken(accessToken)
}
