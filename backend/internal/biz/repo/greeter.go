package repo

import (
	"context"

	"backend/internal/biz/entity"
)

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	Save(context.Context, *entity.Greeter) (*entity.Greeter, error)
	Update(context.Context, *entity.Greeter) (*entity.Greeter, error)
	FindByID(context.Context, int64) (*entity.Greeter, error)
	ListByHello(context.Context, string) ([]*entity.Greeter, error)
	ListAll(context.Context) ([]*entity.Greeter, error)
}

// ArticleRepo is an Article repo.
type ArticleRepo interface {
	ListArticles(context.Context) ([]*entity.Article, error)
}
