package api

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"go-template/internal/core/config"
	"go-template/internal/core/ports"
	"go-template/internal/log"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	cfg                *config.Configuration
	translationService ports.TranslationService
}

func NewServer(cfg *config.Configuration, translationService ports.TranslationService) *Server {
	return &Server{
		cfg:                cfg,
		translationService: translationService,
	}
}

func (s *Server) CreateTranslation(ctx context.Context, request CreateTranslationRequestObject) (CreateTranslationResponseObject, error) {
	log.Info(ctx, "[CreateTranslation] <- Enter", "request", request)
	err := validateCreateTranslationRequest(request.Body)
	if err != nil {
		log.Error(ctx, "[CreateTranslation] <- Leave invalid request", "error", err)
		return CreateTranslation400JSONResponse{N400JSONResponse{Message: err.Error()}}, nil
	}
	req := ports.NewTranslationRequest(request.Body.Key, request.Body.Translation, request.Body.LanguageCode)
	id, err := s.translationService.CreateTranslation(req)
	if err != nil {
		log.Error(ctx, "[CreateTranslation] <- Leave with error", "error", err)
		return CreateTranslation500JSONResponse{}, nil
	}
	log.Info(ctx, "[CreateTranslation] <- Leave")
	return CreateTranslation201JSONResponse{Id: id.String()}, nil
}

func validateCreateTranslationRequest(body *CreateTranslationJSONRequestBody) error {
	if strings.TrimSpace(body.Key) == "" {
		return errors.New("missing or empty key")
	}
	if strings.TrimSpace(body.Translation) == "" {
		return errors.New("missing or empty translation")
	}
	if strings.TrimSpace(body.LanguageCode) == "" {
		return errors.New("missing or empty languageCode")
	}
	return nil
}

func (s *Server) GetAllTranslations(ctx context.Context, request GetAllTranslationsRequestObject) (GetAllTranslationsResponseObject, error) {
	log.Info(ctx, "[GetAllTranslations] <- Enter")
	translations, err := s.translationService.GetAllTranslations(request.Params.LanguageCode)
	if err != nil {
		log.Error(ctx, "[GetAllTranslations] <- Leave with error", "error", err)
		return GetAllTranslations500JSONResponse{}, nil
	}
	log.Info(ctx, "[GetAllTranslations] <- Leave ")
	return GetAllTranslations200JSONResponse(toTranslationResponse(translations)), nil
}

func (s *Server) Check(ctx context.Context, _ CheckRequestObject) (CheckResponseObject, error) {
	log.Info(ctx, "[Check] <- Enter and return directly status 200")
	return Check200Response{}, nil
}

// RegisterStatic add method to the mux that are not documented in the API.
func RegisterStatic(mux *chi.Mux) {
	mux.Get("/", documentation)
	mux.Get("/static/docs/api/api.yaml", swagger)
	mux.Get("/favicon.ico", favicon)
}

func documentation(w http.ResponseWriter, _ *http.Request) {
	writeFile("api/spec.html", "text/html; charset=UTF-8", w)
}

func favicon(w http.ResponseWriter, _ *http.Request) {
	writeFile("api/polygon.png", "image/png", w)
}

func swagger(w http.ResponseWriter, _ *http.Request) {
	writeFile("api/api.yaml", "text/html; charset=UTF-8", w)
}

func writeFile(path string, mimeType string, w http.ResponseWriter) {
	f, err := os.ReadFile(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(f)
}
