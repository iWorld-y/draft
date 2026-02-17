package server

import (
	stdhttp "net/http"
	"strings"
	"time"

	v1 "backend/api/helloworld/v1"
	authctx "backend/internal/auth"
	"backend/internal/conf"
	"backend/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// corsFilter CORS 过滤器
func corsFilter() http.FilterFunc {
	return func(h stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			origin := r.Header.Get("Origin")

			// 允许的源
			allowedOrigins := []string{
				"http://localhost:8123",
				"http://127.0.0.1:8123",
				"http://localhost:3000",
				"http://localhost:5173",
			}

			isAllowed := false
			for _, allowed := range allowedOrigins {
				if origin == allowed || origin == "" {
					isAllowed = true
					break
				}
			}

			if isAllowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			// 处理 OPTIONS 预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(200)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, dictSvc *service.DictionaryService, learnSvc *service.LearningService, authSvc *service.AuthService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			requestLogMiddleware(logger),
		),
		http.Filter(requestContextFilter(authSvc), corsFilter()),
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

	// 注册认证相关路由
	registerAuthRoutes(srv, authSvc)

	// 注册词典相关路由
	registerDictionaryRoutes(srv, dictSvc, authSvc)

	// 注册学习相关路由
	registerLearningRoutes(srv, learnSvc, authSvc)

	return srv
}

func setRefreshCookie(ctx http.Context, refreshToken string) {
	stdhttp.SetCookie(ctx.Response(), &stdhttp.Cookie{
		Name:     service.RefreshCookieName,
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		HttpOnly: true,
		SameSite: stdhttp.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}

func clearRefreshCookie(ctx http.Context) {
	stdhttp.SetCookie(ctx.Response(), &stdhttp.Cookie{
		Name:     service.RefreshCookieName,
		Value:    "",
		Path:     "/api/v1/auth",
		HttpOnly: true,
		SameSite: stdhttp.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

func withAuth(authSvc *service.AuthService, handler func(ctx http.Context, userID int64) error) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		if userID, ok := authctx.UserIDFromContext(ctx); ok && userID > 0 {
			return handler(ctx, userID)
		}

		authz := strings.TrimSpace(ctx.Request().Header.Get("Authorization"))
		if !strings.HasPrefix(authz, "Bearer ") {
			return ctx.JSON(401, map[string]interface{}{"code": 401, "message": "未授权"})
		}

		token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
		userID, err := authSvc.ParseAccessToken(token)
		if err != nil {
			return ctx.JSON(401, map[string]interface{}{"code": 401, "message": "未授权"})
		}

		return handler(ctx, userID)
	}
}

// registerAuthRoutes 注册认证路由
func registerAuthRoutes(srv *http.Server, svc *service.AuthService) {
	route := srv.Route("/api/v1/auth")

	route.POST("/register", func(ctx http.Context) error {
		var req service.AuthRequest
		if err := ctx.Bind(&req); err != nil {
			return err
		}
		resp, refreshToken, err := svc.Register(ctx, &req)
		if err != nil {
			return err
		}
		setRefreshCookie(ctx, refreshToken)
		return ctx.JSON(200, map[string]interface{}{"code": 0, "data": resp})
	})

	route.POST("/login", func(ctx http.Context) error {
		var req service.AuthRequest
		if err := ctx.Bind(&req); err != nil {
			return err
		}
		resp, refreshToken, err := svc.Login(ctx, &req)
		if err != nil {
			return err
		}
		setRefreshCookie(ctx, refreshToken)
		return ctx.JSON(200, map[string]interface{}{"code": 0, "data": resp})
	})

	route.POST("/refresh", func(ctx http.Context) error {
		cookie, err := ctx.Request().Cookie(service.RefreshCookieName)
		if err != nil || cookie == nil || strings.TrimSpace(cookie.Value) == "" {
			return ctx.JSON(401, map[string]interface{}{"code": 401, "message": "未授权"})
		}
		resp, newRefreshToken, err := svc.Refresh(ctx, cookie.Value)
		if err != nil {
			clearRefreshCookie(ctx)
			return ctx.JSON(401, map[string]interface{}{"code": 401, "message": "未授权"})
		}
		setRefreshCookie(ctx, newRefreshToken)
		return ctx.JSON(200, map[string]interface{}{"code": 0, "data": resp})
	})

	route.POST("/logout", func(ctx http.Context) error {
		cookie, err := ctx.Request().Cookie(service.RefreshCookieName)
		if err == nil && cookie != nil {
			_ = svc.Logout(ctx, cookie.Value)
		}
		clearRefreshCookie(ctx)
		return ctx.JSON(200, map[string]interface{}{"code": 0, "data": map[string]interface{}{"success": true}})
	})

	route.GET("/me", withAuth(svc, func(ctx http.Context, userID int64) error {
		user, err := svc.Me(ctx, userID)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{"code": 0, "data": user})
	}))
}

// registerDictionaryRoutes 注册词典路由
func registerDictionaryRoutes(srv *http.Server, svc *service.DictionaryService, authSvc *service.AuthService) {
	route := srv.Route("/api/v1")

	// 创建词典
	route.POST("/dictionaries", withAuth(authSvc, func(ctx http.Context, userID int64) error {
		var req service.CreateDictionaryRequest
		if err := ctx.Bind(&req); err != nil {
			return err
		}
		resp, err := svc.CreateDictionary(authctx.WithUserID(ctx, userID), &req)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	}))

	// 词典列表
	route.GET("/dictionaries", withAuth(authSvc, func(ctx http.Context, userID int64) error {
		resp, err := svc.ListDictionaries(authctx.WithUserID(ctx, userID))
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	}))

	// 上传词典文件
	route.POST("/dictionaries/upload", withAuth(authSvc, func(ctx http.Context, userID int64) error {
		resp, err := svc.UploadDictionary(authctx.WithUserID(ctx, userID))
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	}))

	// 查询上传任务状态
	route.GET("/dictionaries/upload/status/{task_id}", withAuth(authSvc, func(ctx http.Context, userID int64) error {
		req := &service.GetUploadStatusRequest{
			TaskID: ctx.Vars().Get("task_id"),
		}
		resp, err := svc.GetUploadStatus(authctx.WithUserID(ctx, userID), req)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	}))
}

// registerLearningRoutes 注册学习路由
func registerLearningRoutes(srv *http.Server, svc *service.LearningService, authSvc *service.AuthService) {
	route := srv.Route("/api/v1")

	// 获取今日学习任务
	route.GET("/learning/today-tasks", withAuth(authSvc, func(ctx http.Context, userID int64) error {
		req := &service.TodayTasksRequest{
			DictID: ctx.Query().Get("dict_id"),
			Limit:  ctx.Query().Get("limit"),
		}
		resp, err := svc.GetTodayTasks(authctx.WithUserID(ctx, userID), req)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	}))

	// 提交学习结果
	route.POST("/learning/submit", withAuth(authSvc, func(ctx http.Context, userID int64) error {
		var req service.SubmitLearningRequest
		if err := ctx.Bind(&req); err != nil {
			return err
		}
		resp, err := svc.SubmitLearning(authctx.WithUserID(ctx, userID), &req)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	}))
}
