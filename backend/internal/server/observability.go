package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	stdhttp "net/http"
	"strings"
	"time"

	authctx "backend/internal/auth"
	"backend/internal/service"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

const requestIDHeader = "X-Request-Id"

func requestContextFilter(authSvc *service.AuthService) http.FilterFunc {
	return func(h stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			requestID := strings.TrimSpace(r.Header.Get(requestIDHeader))
			if requestID == "" {
				requestID = newRequestID()
			}

			ctx := authctx.WithRequestID(r.Context(), requestID)
			if userID := parseUserIDFromHeader(authSvc, r.Header.Get("Authorization")); userID > 0 {
				ctx = authctx.WithUserID(ctx, userID)
			}

			w.Header().Set(requestIDHeader, requestID)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func requestLogMiddleware(logger log.Logger) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			start := time.Now()
			requestID, _ := authctx.RequestIDFromContext(ctx)
			if requestID == "" {
				requestID = newRequestID()
				ctx = authctx.WithRequestID(ctx, requestID)
			}
			userID, _ := authctx.UserIDFromContext(ctx)

			kind := "unknown"
			operation := "unknown"
			if tr, ok := transport.FromServerContext(ctx); ok {
				kind = string(tr.Kind())
				operation = tr.Operation()
			}

			l := log.With(
				logger,
				"kind", kind,
				"operation", operation,
				"request_id", requestID,
				"user_id", userID,
			)
			h := log.NewHelper(l)

			reply, err := next(ctx, req)
			latencyMS := time.Since(start).Milliseconds()
			if err != nil {
				se := kerrors.FromError(err)
				h.Errorf(
					"request failed latency_ms=%d error_code=%d error_reason=%s error_message=%v",
					latencyMS,
					se.Code,
					se.Reason,
					err,
				)
				return reply, err
			}

			h.Infof("request handled latency_ms=%d", latencyMS)
			return reply, nil
		}
	}
}

func parseUserIDFromHeader(authSvc *service.AuthService, authz string) int64 {
	if authSvc == nil {
		return 0
	}
	token := strings.TrimSpace(authz)
	if !strings.HasPrefix(token, "Bearer ") {
		return 0
	}
	userID, err := authSvc.ParseAccessToken(strings.TrimSpace(strings.TrimPrefix(token, "Bearer ")))
	if err != nil || userID <= 0 {
		return 0
	}
	return userID
}

func newRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b[:])
}
