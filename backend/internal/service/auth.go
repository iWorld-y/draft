package service

import (
	"context"

	v1 "backend/api/helloworld/v1"
	authctx "backend/internal/auth"
	"backend/internal/biz"
)

type AuthService struct {
	v1.UnimplementedAuthServer

	uc *biz.AuthUseCase
}

func NewAuthService(uc *biz.AuthUseCase) *AuthService {
	return &AuthService{uc: uc}
}

func (s *AuthService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.AuthReply, error) {
	result, err := s.uc.Register(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &v1.AuthReply{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User: &v1.AuthUser{
			Id:       result.User.ID,
			Username: result.User.Username,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.AuthReply, error) {
	result, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &v1.AuthReply{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User: &v1.AuthUser{
			Id:       result.User.ID,
			Username: result.User.Username,
		},
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, req *v1.RefreshRequest) (*v1.AuthReply, error) {
	result, err := s.uc.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &v1.AuthReply{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User: &v1.AuthUser{
			Id:       result.User.ID,
			Username: result.User.Username,
		},
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutReply, error) {
	if err := s.uc.Logout(ctx, req.RefreshToken); err != nil {
		return nil, err
	}
	return &v1.LogoutReply{Success: true}, nil
}

func (s *AuthService) Me(ctx context.Context, _ *v1.MeRequest) (*v1.MeReply, error) {
	userID, ok := authctx.UserIDFromContext(ctx)
	if !ok || userID <= 0 {
		return nil, biz.ErrUnauthorized
	}
	user, err := s.uc.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &v1.MeReply{User: &v1.AuthUser{Id: user.ID, Username: user.Username}}, nil
}

func (s *AuthService) ParseAccessToken(accessToken string) (int64, error) {
	return s.uc.ParseAccessToken(accessToken)
}
