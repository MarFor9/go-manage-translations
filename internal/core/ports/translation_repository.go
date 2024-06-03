package ports

import (
	"github.com/google/uuid"
	"go-template/internal/core/domain"
)

type TranslationRepository interface {
	GetAllTranslations(languageCode string) ([]domain.Translation, error)
	CreateTranslation(request *TranslationRequest) (uuid.UUID, error)
}
