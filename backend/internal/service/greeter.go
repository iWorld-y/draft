package service

import (
	"context"

	v1 "backend/api/helloworld/v1"
	"backend/internal/biz"
	"backend/internal/biz/entity"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &entity.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

// ListArticles implements helloworld.GreeterServer.
func (s *GreeterService) ListArticles(ctx context.Context, in *v1.ListArticlesRequest) (*v1.ListArticlesReply, error) {
	articles, err := s.uc.ListArticles(ctx)
	if err != nil {
		return nil, err
	}
	reply := &v1.ListArticlesReply{}
	for _, a := range articles {
		reply.Articles = append(reply.Articles, &v1.Article{
			Title:   a.Title,
			Content: a.Content,
		})
	}
	return reply, nil
}
