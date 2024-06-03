// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

const (
	ApiKeyAuthScopes = "ApiKeyAuth.Scopes"
)

// CreateTranslationRequest defines model for CreateTranslationRequest.
type CreateTranslationRequest struct {
	Key          string `json:"key"`
	LanguageCode string `json:"languageCode"`
	Translation  string `json:"translation"`
}

// GenericErrorMessage defines model for GenericErrorMessage.
type GenericErrorMessage struct {
	Message string `json:"message"`
}

// UUIDResponse defines model for UUIDResponse.
type UUIDResponse struct {
	Id string `json:"id"`
}

// N400 defines model for 400.
type N400 = GenericErrorMessage

// N401 defines model for 401.
type N401 = GenericErrorMessage

// N500 defines model for 500.
type N500 = GenericErrorMessage

// GetAllTranslationsParams defines parameters for GetAllTranslations.
type GetAllTranslationsParams struct {
	// LanguageCode Language code to filter
	LanguageCode string `form:"language_code" json:"language_code"`
}

// CreateTranslationJSONRequestBody defines body for CreateTranslation for application/json ContentType.
type CreateTranslationJSONRequestBody = CreateTranslationRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Check if the service is up
	// (GET /health)
	Check(w http.ResponseWriter, r *http.Request)
	// Get all translations
	// (GET /v1/translations)
	GetAllTranslations(w http.ResponseWriter, r *http.Request, params GetAllTranslationsParams)
	// Create a new translation
	// (POST /v1/translations)
	CreateTranslation(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Check if the service is up
// (GET /health)
func (_ Unimplemented) Check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get all translations
// (GET /v1/translations)
func (_ Unimplemented) GetAllTranslations(w http.ResponseWriter, r *http.Request, params GetAllTranslationsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Create a new translation
// (POST /v1/translations)
func (_ Unimplemented) CreateTranslation(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Check operation middleware
func (siw *ServerInterfaceWrapper) Check(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Check(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetAllTranslations operation middleware
func (siw *ServerInterfaceWrapper) GetAllTranslations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, ApiKeyAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAllTranslationsParams

	// ------------- Required query parameter "language_code" -------------

	if paramValue := r.URL.Query().Get("language_code"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "language_code"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "language_code", r.URL.Query(), &params.LanguageCode)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "language_code", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAllTranslations(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateTranslation operation middleware
func (siw *ServerInterfaceWrapper) CreateTranslation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, ApiKeyAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateTranslation(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/health", wrapper.Check)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/translations", wrapper.GetAllTranslations)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/translations", wrapper.CreateTranslation)
	})

	return r
}

type N400JSONResponse GenericErrorMessage

type N401JSONResponse GenericErrorMessage

type N500JSONResponse GenericErrorMessage

type CheckRequestObject struct {
}

type CheckResponseObject interface {
	VisitCheckResponse(w http.ResponseWriter) error
}

type Check200Response struct {
}

func (response Check200Response) VisitCheckResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type GetAllTranslationsRequestObject struct {
	Params GetAllTranslationsParams
}

type GetAllTranslationsResponseObject interface {
	VisitGetAllTranslationsResponse(w http.ResponseWriter) error
}

type GetAllTranslations200JSONResponse map[string]string

func (response GetAllTranslations200JSONResponse) VisitGetAllTranslationsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetAllTranslations400JSONResponse struct{ N400JSONResponse }

func (response GetAllTranslations400JSONResponse) VisitGetAllTranslationsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetAllTranslations401JSONResponse struct{ N401JSONResponse }

func (response GetAllTranslations401JSONResponse) VisitGetAllTranslationsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetAllTranslations500JSONResponse struct{ N500JSONResponse }

func (response GetAllTranslations500JSONResponse) VisitGetAllTranslationsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type CreateTranslationRequestObject struct {
	Body *CreateTranslationJSONRequestBody
}

type CreateTranslationResponseObject interface {
	VisitCreateTranslationResponse(w http.ResponseWriter) error
}

type CreateTranslation201JSONResponse UUIDResponse

func (response CreateTranslation201JSONResponse) VisitCreateTranslationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type CreateTranslation400JSONResponse struct{ N400JSONResponse }

func (response CreateTranslation400JSONResponse) VisitCreateTranslationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type CreateTranslation401JSONResponse struct{ N401JSONResponse }

func (response CreateTranslation401JSONResponse) VisitCreateTranslationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type CreateTranslation500JSONResponse struct{ N500JSONResponse }

func (response CreateTranslation500JSONResponse) VisitCreateTranslationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Check if the service is up
	// (GET /health)
	Check(ctx context.Context, request CheckRequestObject) (CheckResponseObject, error)
	// Get all translations
	// (GET /v1/translations)
	GetAllTranslations(ctx context.Context, request GetAllTranslationsRequestObject) (GetAllTranslationsResponseObject, error)
	// Create a new translation
	// (POST /v1/translations)
	CreateTranslation(ctx context.Context, request CreateTranslationRequestObject) (CreateTranslationResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHttpHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHttpMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// Check operation middleware
func (sh *strictHandler) Check(w http.ResponseWriter, r *http.Request) {
	var request CheckRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Check(ctx, request.(CheckRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Check")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CheckResponseObject); ok {
		if err := validResponse.VisitCheckResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetAllTranslations operation middleware
func (sh *strictHandler) GetAllTranslations(w http.ResponseWriter, r *http.Request, params GetAllTranslationsParams) {
	var request GetAllTranslationsRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetAllTranslations(ctx, request.(GetAllTranslationsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetAllTranslations")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetAllTranslationsResponseObject); ok {
		if err := validResponse.VisitGetAllTranslationsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// CreateTranslation operation middleware
func (sh *strictHandler) CreateTranslation(w http.ResponseWriter, r *http.Request) {
	var request CreateTranslationRequestObject

	var body CreateTranslationJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateTranslation(ctx, request.(CreateTranslationRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateTranslation")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateTranslationResponseObject); ok {
		if err := validResponse.VisitCreateTranslationResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}