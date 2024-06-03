package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go-template/internal/api"
	"go-template/internal/core/config"
	"go-template/internal/core/services"
	"go-template/internal/db"
	"go-template/internal/errors"
	"go-template/internal/log"
	"go-template/internal/repositories"
	"net/http"
	"os"
)

func main() {
	log.Info(context.Background(), "Starting server")
	cfg, err := config.Load()
	if err != nil {
		log.Error(context.Background(), "Error loading configuration: %s", err)
		return
	}

	ctx, cancel := context.WithCancel(log.NewContext(context.Background(), cfg.Log.Level, cfg.Log.Mode, os.Stdout))
	defer cancel()

	storage, err := db.NewStorage(cfg.Database.URL)
	if err != nil {
		log.Error(ctx, "Error connecting to database", "err", err)
		return
	}

	//repositories
	translationRepository := repositories.NewTranslation(*storage)

	//services
	translationService := services.NewTranslation(translationRepository)

	mux := chi.NewRouter()
	mux.Use(
		chiMiddleware.RequestID,
		log.ChiMiddleware(ctx),
		chiMiddleware.Recoverer,
		cors.AllowAll().Handler,
		chiMiddleware.NoCache,
	)
	api.HandlerWithOptions(
		api.NewStrictHandlerWithOptions(
			api.NewServer(cfg, translationService),
			middlewares(ctx, cfg),
			api.StrictHTTPServerOptions{
				RequestErrorHandlerFunc:  errors.RequestErrorHandlerFunc,
				ResponseErrorHandlerFunc: errors.ResponseErrorHandlerFunc,
			}),
		api.ChiServerOptions{
			BaseRouter:       mux,
			ErrorHandlerFunc: errorHandlerFunc,
		},
	)

	api.RegisterStatic(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: mux,
	}

	log.Info(ctx, fmt.Sprintf("Starting server on %s:%d", cfg.ServerUrl, cfg.ServerPort))
	if err := server.ListenAndServe(); err != nil {
		log.Error(ctx, "starting HTTP UI API server", "err", err)
	}
}

func middlewares(ctx context.Context, cfg *config.Configuration) []api.StrictMiddlewareFunc {
	return []api.StrictMiddlewareFunc{
		api.LogMiddleware(ctx),
		api.ApiKeyAuthMiddleware(ctx, cfg),
	}
}

func errorHandlerFunc(w http.ResponseWriter, _ *http.Request, err error) {
	switch err.(type) {
	case *api.InvalidParamFormatError:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": err.Error()})
	default:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
