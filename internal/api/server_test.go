package api

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-template/internal/api/tests"
	"go-template/internal/core/ports"
	"go-template/internal/core/services"
	"go-template/internal/log"
	"go-template/internal/repositories"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServer_Check(t *testing.T) {
	server := NewServer(nil, nil)
	ctx := log.NewContext(context.Background(), log.LevelDebug, log.OutputText, os.Stdout)

	handler := getHandler(ctx, server)

	type testConfig struct {
		name     string
		httpCode int
	}

	for _, tc := range []testConfig{
		{
			name:     "Health check",
			httpCode: 200,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			require.NoError(t, err)

			handler.ServeHTTP(rr, req)
			require.Equal(t, tc.httpCode, rr.Code)
		})
	}
}
func TestServer_GetAllTranslations(t *testing.T) {
	cfg.TranslationAPIKey = "test"

	translationRepository := repositories.NewTranslation(*storage)
	translationService := services.NewTranslation(translationRepository)

	server := NewServer(&cfg, translationService)
	ctx := log.NewContext(context.Background(), log.LevelDebug, log.OutputText, os.Stdout)

	handler := getHandler(ctx, server)

	//init database
	_, _ = translationRepository.CreateTranslation(ports.NewTranslationRequest("input_avatarLabel", "Upload Avatar", "en"))
	_, _ = translationRepository.CreateTranslation(ports.NewTranslationRequest("input_avatarHelperText", "Drag and Drop or", "en"))
	_, _ = translationRepository.CreateTranslation(ports.NewTranslationRequest("input_firstNameLabel", "First name", "en"))

	defer storage.GormDB.Exec("DELETE FROM translations")

	type expected struct {
		httpCode int
		response GetAllTranslationsResponseObject
	}

	type testConfig struct {
		name     string
		apiKey   string
		expected expected
		url      string
	}

	for _, tc := range []testConfig{
		{
			name:   "Happy path - get all translations by language code",
			apiKey: cfg.TranslationAPIKey,
			expected: expected{
				httpCode: 200,
				response: GetAllTranslations200JSONResponse{
					"input_avatarLabel":      "Upload Avatar",
					"input_avatarHelperText": "Drag and Drop or",
					"input_firstNameLabel":   "First name",
				},
			},
			url: "/v1/translations?language_code=en",
		},
		{
			name:   "Translation not found",
			apiKey: cfg.TranslationAPIKey,
			expected: expected{
				httpCode: 200,
				response: GetAllTranslations200JSONResponse{},
			},
			url: "/v1/translations?language_code=unknown",
		},
		{
			name:   "Unauthorized",
			apiKey: "wrong-api-key",
			expected: expected{
				httpCode: 401,
				response: GetAllTranslations401JSONResponse{},
			},
			url: "/v1/translations?language_code=en",
		},
		{
			name:   "Bad request - empty language code",
			apiKey: cfg.TranslationAPIKey,
			expected: expected{
				httpCode: 400,
				response: GetAllTranslations400JSONResponse{},
			},
			url: "/v1/translations",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
			req.Header.Set("X-API-Key", tc.apiKey)
			require.NoError(t, err)

			handler.ServeHTTP(rr, req)
			require.Equal(t, tc.expected.httpCode, rr.Code)
			switch tc.expected.httpCode {
			case http.StatusOK:
				var response GetAllTranslations200JSONResponse
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
				assert.Equal(t, tc.expected.response, response)
			case http.StatusBadRequest:
				assert.Contains(t, rr.Body.String(), "Query argument language_code is required, but not found")
			}
		})
	}
}

func TestServer_CreateTranslation(t *testing.T) {
	cfg.TranslationAPIKey = "test"

	translationRepository := repositories.NewTranslation(*storage)
	translationService := services.NewTranslation(translationRepository)

	server := NewServer(&cfg, translationService)
	ctx := log.NewContext(context.Background(), log.LevelDebug, log.OutputText, os.Stdout)

	handler := getHandler(ctx, server)

	type expected struct {
		httpCode int
		response CreateTranslationResponseObject
		errorMsg string
	}

	type testConfig struct {
		name     string
		apiKey   string
		body     CreateTranslationJSONRequestBody
		expected expected
	}

	for _, tc := range []testConfig{
		{
			name:   "Happy path - translation created",
			apiKey: cfg.TranslationAPIKey,
			body: CreateTranslationJSONRequestBody{
				Key:          "input_lastNameLabel",
				Translation:  "Last name",
				LanguageCode: "en",
			},
			expected: expected{
				httpCode: 201,
				response: CreateTranslation201JSONResponse{},
			},
		},
		{
			name:   "Unauthorized",
			apiKey: "wrong-api-key",
			body: CreateTranslationJSONRequestBody{
				Key:          "input_lastNameLabel",
				Translation:  "Last name",
				LanguageCode: "en",
			},
			expected: expected{
				httpCode: 401,
				response: CreateTranslation401JSONResponse{},
			},
		},
		{
			name:   "Bad request - empty languageCode",
			apiKey: cfg.TranslationAPIKey,
			body: CreateTranslationJSONRequestBody{
				Key:          "input_lastNameLabel",
				Translation:  "Last name",
				LanguageCode: "",
			},
			expected: expected{
				httpCode: 400,
				response: CreateTranslation400JSONResponse{},
				errorMsg: "missing or empty languageCode",
			},
		},
		{
			name:   "Bad request - empty key",
			apiKey: cfg.TranslationAPIKey,
			body: CreateTranslationJSONRequestBody{
				Key:          "",
				Translation:  "Last name",
				LanguageCode: "en",
			},
			expected: expected{
				httpCode: 400,
				response: CreateTranslation400JSONResponse{},
				errorMsg: "missing or empty key",
			},
		},
		{
			name:   "Bad request - empty translation",
			apiKey: cfg.TranslationAPIKey,
			body: CreateTranslationJSONRequestBody{
				Key:          "input_lastNameLabel",
				Translation:  "",
				LanguageCode: "en",
			},
			expected: expected{
				httpCode: 400,
				response: CreateTranslation400JSONResponse{},
				errorMsg: "missing or empty translation",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/v1/translations", tests.JSONBody(t, tc.body))
			req.Header.Set("X-API-Key", tc.apiKey)
			require.NoError(t, err)

			handler.ServeHTTP(rr, req)
			require.Equal(t, tc.expected.httpCode, rr.Code)
			switch tc.expected.httpCode {
			case http.StatusCreated:
				var response CreateTranslation201JSONResponse
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
				getAllTranslations, err := translationRepository.GetAllTranslations("en")
				require.NoError(t, err)
				assert.Equal(t, getAllTranslations[0].Key, tc.body.Key)
				assert.Equal(t, getAllTranslations[0].Translation, tc.body.Translation)
				assert.Equal(t, getAllTranslations[0].LanguageCode, tc.body.LanguageCode)
			case http.StatusBadRequest:
				var response CreateTranslation400JSONResponse
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
				assert.Equal(t, tc.expected.errorMsg, response.Message)
			}
		})
	}
}
