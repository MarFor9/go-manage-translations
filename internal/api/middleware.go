package api

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"go-template/internal/core/config"
	apiErrors "go-template/internal/errors"
	"go-template/internal/log"
	"net/http"
)

// LogMiddleware returns a middleware that adds general log configuration to each context request
func LogMiddleware(ctx context.Context) StrictMiddlewareFunc {
	return func(f StrictHandlerFunc, operationID string) StrictHandlerFunc {
		return func(ctxReq context.Context, w http.ResponseWriter, r *http.Request, args interface{}) (interface{}, error) {
			ctx := log.CopyFromContext(ctx, ctxReq)
			if reqID := middleware.GetReqID(ctxReq); reqID != "" {
				ctx = log.With(ctx, "req-id", reqID)
			}
			return f(ctx, w, r, args)
		}
	}
}

// ApiKeyAuthMiddleware returns a middleware that performs an http api key authorization for endpoints configured with
// ApiKeyAuth auth in the api spec.
// In uses the ApiKeyAuthScopes value in context to figure if and endpoint needs authorization or not, because this
// value is injected automatically by openapi when basic auth is selected
func ApiKeyAuthMiddleware(ctx context.Context, cfg *config.Configuration) StrictMiddlewareFunc {
	return func(f StrictHandlerFunc, operationID string) StrictHandlerFunc {
		return func(ctxReq context.Context, w http.ResponseWriter, r *http.Request, args interface{}) (interface{}, error) {
			log.Info(ctx, "ApiKeyAuthMiddleware", "authScopes", ctxReq.Value(ApiKeyAuthScopes))

			if ctxReq.Value(ApiKeyAuthScopes) != nil {
				apikey := cfg.TranslationAPIKey
				getApiKey := r.Header.Get("X-API-KEY")
				if getApiKey == "" {
					return nil, apiErrors.AuthError{Err: errors.New("api key not provided. Unauthorized")}
				}
				if getApiKey != apikey {
					return nil, apiErrors.AuthError{Err: errors.New("api key not valid. Unauthorized")}
				}
				return f(ctx, w, r, args)
			}

			return f(ctx, w, r, args)
		}
	}
}
