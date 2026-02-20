package server

import (
	"context"
	stdhttp "net/http"
	"strings"

	v1 "backend/api/helloworld/v1"
	authctx "backend/internal/auth"
	"backend/internal/conf"
	"backend/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// corsFilter CORS 过滤器
func corsFilter() http.FilterFunc {
	return func(h stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			origin := r.Header.Get("Origin")

			// Allow any origin for now to solve connectivity issues
			// In production, this should be restricted to specific domains

			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(200)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func authContextMiddleware(authSvc *service.AuthService) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			authz := strings.TrimSpace(tr.RequestHeader().Get("Authorization"))
			if !strings.HasPrefix(authz, "Bearer ") {
				return handler(ctx, req)
			}

			token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
			userID, err := authSvc.ParseAccessToken(token)
			if err != nil {
				return nil, err
			}

			return handler(authctx.WithUserID(ctx, userID), req)
		}
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, dictSvc *service.DictionaryService, learnSvc *service.LearningService, authSvc *service.AuthService, logger log.Logger) *http.Server {
	_ = logger

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			authContextMiddleware(authSvc),
		),
		http.Filter(corsFilter()),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	v1.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterAuthHTTPServer(srv, authSvc)
	v1.RegisterDictionaryHTTPServer(srv, dictSvc)
	v1.RegisterLearningHTTPServer(srv, learnSvc)

	return srv
}
