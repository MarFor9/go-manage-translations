package services

import (
	"github.com/google/uuid"
	"go-template/internal/core/domain"
	"go-template/internal/core/ports"
)

type Translation struct {
	repo ports.TranslationRepository
}

func NewTranslation(repo ports.TranslationRepository) *Translation {
	return &Translation{
		repo: repo,
	}
}

func (t *Translation) GetAllTranslations(languageCode string) ([]domain.Translation, error) {
	return t.repo.GetAllTranslations(languageCode)
}

func (t *Translation) CreateTranslation(translationRequest *ports.TranslationRequest) (uuid.UUID, error) {
	return t.repo.CreateTranslation(translationRequest)
}
