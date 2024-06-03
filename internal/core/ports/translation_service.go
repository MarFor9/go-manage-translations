package ports

import (
	"github.com/google/uuid"
	"go-template/internal/core/domain"
)

type TranslationService interface {
	GetAllTranslations(languageCode string) ([]domain.Translation, error)
	CreateTranslation(translationRequest *TranslationRequest) (uuid.UUID, error)
}

type TranslationRequest struct {
	Key          string
	Translation  string
	LanguageCode string
}

func NewTranslationRequest(key, translation, code string) *TranslationRequest {
	return &TranslationRequest{
		Key:          key,
		Translation:  translation,
		LanguageCode: code,
	}
}
