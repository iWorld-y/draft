package server

import (
	v1 "backend/api/helloworld/v1"
	"backend/internal/conf"
	"backend/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, dictSvc *service.DictionaryService, learnSvc *service.LearningService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
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

	// 注册词典相关路由
	registerDictionaryRoutes(srv, dictSvc)

	// 注册学习相关路由
	registerLearningRoutes(srv, learnSvc)

	return srv
}

// registerDictionaryRoutes 注册词典路由
func registerDictionaryRoutes(srv *http.Server, svc *service.DictionaryService) {
	route := srv.Route("/api/v1")

	// 创建词典
	route.POST("/dictionaries", func(ctx http.Context) error {
		var req service.CreateDictionaryRequest
		if err := ctx.Bind(&req); err != nil {
			return err
		}
		resp, err := svc.CreateDictionary(ctx, &req)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	})

	// 词典列表
	route.GET("/dictionaries", func(ctx http.Context) error {
		resp, err := svc.ListDictionaries(ctx)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	})

	// 上传词典文件
	route.POST("/dictionaries/upload", func(ctx http.Context) error {
		resp, err := svc.UploadDictionary(ctx)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	})

	// 查询上传任务状态
	route.GET("/dictionaries/upload/status/{task_id}", func(ctx http.Context) error {
		req := &service.GetUploadStatusRequest{
			TaskID: ctx.Vars().Get("task_id"),
		}
		resp, err := svc.GetUploadStatus(ctx, req)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	})
}

// registerLearningRoutes 注册学习路由
func registerLearningRoutes(srv *http.Server, svc *service.LearningService) {
	route := srv.Route("/api/v1")

	// 获取今日学习任务
	route.GET("/learning/today-tasks", func(ctx http.Context) error {
		req := &service.TodayTasksRequest{
			DictID: ctx.Query().Get("dict_id"),
			Limit:  ctx.Query().Get("limit"),
		}
		resp, err := svc.GetTodayTasks(ctx, req)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	})

	// 提交学习结果
	route.POST("/learning/submit", func(ctx http.Context) error {
		var req service.SubmitLearningRequest
		if err := ctx.Bind(&req); err != nil {
			return err
		}
		resp, err := svc.SubmitLearning(ctx, &req)
		if err != nil {
			return err
		}
		return ctx.JSON(200, map[string]interface{}{
			"code": 0,
			"data": resp,
		})
	})
}
