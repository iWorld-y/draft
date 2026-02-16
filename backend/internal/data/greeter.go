package data

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"backend/internal/biz/entity"
	"backend/internal/biz/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) repo.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// NewArticleRepo .
func NewArticleRepo(data *Data, logger log.Logger) repo.ArticleRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) Save(ctx context.Context, g *entity.Greeter) (*entity.Greeter, error) {
	return g, nil
}

func (r *greeterRepo) Update(ctx context.Context, g *entity.Greeter) (*entity.Greeter, error) {
	return g, nil
}

func (r *greeterRepo) FindByID(context.Context, int64) (*entity.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListByHello(context.Context, string) ([]*entity.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListAll(context.Context) ([]*entity.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListArticles(ctx context.Context) ([]*entity.Article, error) {
	dir := "../../docs/"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var articles []*entity.Article
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				continue
			}
			articles = append(articles, &entity.Article{
				Title:   entry.Name(),
				Content: string(content),
			})
		}
	}
	return articles, nil
}
