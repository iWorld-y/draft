package service

import (
	"context"

	"backend/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

const RefreshCookieName = "refresh_token"

type AuthService struct {
	uc  *biz.AuthUseCase
	log *log.Helper
}

func NewAuthService(uc *biz.AuthUseCase, logger log.Logger) *AuthService {
	return &AuthService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
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
	s.log.WithContext(ctx).Infof("Register req: %+v", req)
	result, err := s.uc.Register(ctx, req.Username, req.Password)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Register failed req=%+v err=%v", req, err)
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
	s.log.WithContext(ctx).Infof("Login req: %+v", req)
	result, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Login failed req=%+v err=%v", req, err)
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
	s.log.WithContext(ctx).Infof("Refresh req: refreshToken=%s", refreshToken)
	result, err := s.uc.Refresh(ctx, refreshToken)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Refresh failed refreshToken=%s err=%v", refreshToken, err)
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
	s.log.WithContext(ctx).Infof("Logout req: refreshToken=%s", refreshToken)
	err := s.uc.Logout(ctx, refreshToken)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Logout failed refreshToken=%s err=%v", refreshToken, err)
		return err
	}
	return nil
}

func (s *AuthService) Me(ctx context.Context, userID int64) (*AuthUser, error) {
	s.log.WithContext(ctx).Infof("Me req: userID=%d", userID)
	user, err := s.uc.GetUserByID(ctx, userID)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Me failed userID=%d err=%v", userID, err)
		return nil, err
	}
	return &AuthUser{ID: user.ID, Username: user.Username}, nil
}

func (s *AuthService) ParseAccessToken(accessToken string) (int64, error) {
	s.log.Infof("ParseAccessToken req: accessToken=%s", accessToken)
	userID, err := s.uc.ParseAccessToken(accessToken)
	if err != nil {
		s.log.Errorf("ParseAccessToken failed accessToken=%s err=%v", accessToken, err)
		return 0, err
	}
	return userID, nil
}
