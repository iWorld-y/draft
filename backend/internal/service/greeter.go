package service

import (
	"context"

	v1 "backend/api/helloworld/v1"
	"backend/internal/biz"
	"backend/internal/biz/entity"

	"github.com/go-kratos/kratos/v2/log"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc  *biz.GreeterUsecase
	log *log.Helper
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	return &GreeterService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	s.log.WithContext(ctx).Infof("SayHello req: %+v", in)
	g, err := s.uc.CreateGreeter(ctx, &entity.Greeter{Hello: in.Name})
	if err != nil {
		s.log.WithContext(ctx).Errorf("SayHello failed req=%+v err=%v", in, err)
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

// ListArticles implements helloworld.GreeterServer.
func (s *GreeterService) ListArticles(ctx context.Context, in *v1.ListArticlesRequest) (*v1.ListArticlesReply, error) {
	s.log.WithContext(ctx).Infof("ListArticles req: %+v", in)
	articles, err := s.uc.ListArticles(ctx)
	if err != nil {
		s.log.WithContext(ctx).Errorf("ListArticles failed req=%+v err=%v", in, err)
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
