package repo

import (
	"context"

	"backend/internal/biz/entity"
)

type UserRepo interface {
	Create(ctx context.Context, user *entity.User) error
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
}

type RefreshTokenRepo interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	GetValidByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	RevokeByHash(ctx context.Context, tokenHash string) error
}
