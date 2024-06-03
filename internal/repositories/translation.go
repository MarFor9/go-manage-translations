package repositories

import (
	"github.com/google/uuid"
	"go-template/internal/core/domain"
	"go-template/internal/core/ports"
	"go-template/internal/db"
)

type Translation struct {
	conn db.Storage
}

func NewTranslation(conn db.Storage) *Translation {
	return &Translation{
		conn: conn,
	}
}

func (t *Translation) GetAllTranslations(languageCode string) ([]domain.Translation, error) {
	var translation []domain.Translation
	result := t.conn.GormDB.Model(&Translation{}).
		Where("language_code = ?", languageCode).
		Find(&translation)
	if result.Error != nil {
		return nil, result.Error
	}
	return translation, nil
}
func (t *Translation) CreateTranslation(request *ports.TranslationRequest) (uuid.UUID, error) {
	translation := domain.Translation{
		ID:           uuid.New(),
		Key:          request.Key,
		Translation:  request.Translation,
		LanguageCode: request.LanguageCode,
	}
	result := t.conn.GormDB.Create(&translation)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return translation.ID, nil
}
