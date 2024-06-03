package api

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go-template/internal/core/config"
	"go-template/internal/db"
	"go-template/internal/db/tests"
	"go-template/internal/errors"
	"go-template/internal/log"
	"net/http"
	"os"
	"testing"
)

var (
	storage *db.Storage
	cfg     config.Configuration
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	conn := lookupPostgresURL()
	if conn == "" {
		conn = "postgres://postgres:postgres@localhost:5435"
	}

	cfgForTesting := config.Configuration{
		Database: config.Database{
			URL: conn,
		},
	}
	s, teardown, err := tests.NewTestStorage(&cfgForTesting)
	defer teardown()
	if err != nil {
		log.Error(ctx, "failed to acquire test database", "err", err)
		os.Exit(1)
	}
	storage = s

	cfg.ServerUrl = "https://testing.env/"

	m.Run()
}

func lookupPostgresURL() string {
	con, ok := os.LookupEnv("POSTGRES_TEST_DATABASE")
	if !ok {
		return ""
	}
	return con
}

func getHandler(ctx context.Context, server *Server) http.Handler {
	mux := chi.NewRouter()
	return HandlerWithOptions(
		NewStrictHandlerWithOptions(
			server,
			middlewares(ctx),
			StrictHTTPServerOptions{
				RequestErrorHandlerFunc:  errors.RequestErrorHandlerFunc,
				ResponseErrorHandlerFunc: errors.ResponseErrorHandlerFunc,
			}),
		ChiServerOptions{
			BaseRouter:       mux,
			ErrorHandlerFunc: errorHandlerFunc,
		},
	)
}

func middlewares(ctx context.Context) []StrictMiddlewareFunc {
	return []StrictMiddlewareFunc{
		LogMiddleware(ctx),
		ApiKeyAuthMiddleware(ctx, &cfg),
	}
}

func errorHandlerFunc(w http.ResponseWriter, _ *http.Request, err error) {
	switch err.(type) {
	case *InvalidParamFormatError:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": err.Error()})
	default:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
