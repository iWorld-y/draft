package biz

import (
	"context"

	v1 "backend/api/helloworld/v1"
	"backend/internal/biz/entity"
	"backend/internal/biz/repo"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo        repo.GreeterRepo
	articleRepo repo.ArticleRepo
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo repo.GreeterRepo, articleRepo repo.ArticleRepo) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, articleRepo: articleRepo}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *entity.Greeter) (*entity.Greeter, error) {
	log.Infof("CreateGreeter: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}

// ListArticles lists all articles in the docs directory.
func (uc *GreeterUsecase) ListArticles(ctx context.Context) ([]*entity.Article, error) {
	return uc.articleRepo.ListArticles(ctx)
}
