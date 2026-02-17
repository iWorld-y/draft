package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"backend/internal/biz/entity"
	"backend/internal/biz/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) repo.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (username, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	if user.Status == 0 {
		user.Status = 1
	}

	if err := r.data.db.QueryRowContext(ctx, query,
		user.Username,
		user.PasswordHash,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID); err != nil {
		r.log.Errorf("failed to create user: %v", err)
		return err
	}
	return nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT id, username, password_hash, status, created_at, updated_at, deleted_at
		FROM users
		WHERE username = $1
	`
	user := &entity.User{}
	err := r.data.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.Errorf("failed to get user by username: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT id, username, password_hash, status, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1
	`
	user := &entity.User{}
	err := r.data.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.Errorf("failed to get user by id: %v", err)
		return nil, err
	}
	return user, nil
}

type refreshTokenRepo struct {
	data *Data
	log  *log.Helper
}

func NewRefreshTokenRepo(data *Data, logger log.Logger) repo.RefreshTokenRepo {
	return &refreshTokenRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *refreshTokenRepo) Create(ctx context.Context, token *entity.RefreshToken) error {
	query := `
		INSERT INTO auth_refresh_tokens (user_id, token_hash, expires_at, revoked_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	token.CreatedAt = time.Now()
	if err := r.data.db.QueryRowContext(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.RevokedAt,
		token.CreatedAt,
	).Scan(&token.ID); err != nil {
		r.log.Errorf("failed to create refresh token: %v", err)
		return err
	}
	return nil
}

func (r *refreshTokenRepo) GetValidByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM auth_refresh_tokens
		WHERE token_hash = $1
		AND revoked_at IS NULL
		AND expires_at > CURRENT_TIMESTAMP
		ORDER BY id DESC
		LIMIT 1
	`
	token := &entity.RefreshToken{}
	err := r.data.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.Errorf("failed to get refresh token: %v", err)
		return nil, err
	}
	return token, nil
}

func (r *refreshTokenRepo) RevokeByHash(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE auth_refresh_tokens
		SET revoked_at = $1
		WHERE token_hash = $2
		AND revoked_at IS NULL
	`
	_, err := r.data.db.ExecContext(ctx, query, time.Now(), tokenHash)
	if err != nil {
		r.log.Errorf("failed to revoke refresh token: %v", err)
		return err
	}
	return nil
}
