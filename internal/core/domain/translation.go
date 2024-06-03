package domain

import "github.com/google/uuid"

type Translation struct {
	ID           uuid.UUID
	Key          string
	Translation  string
	LanguageCode string
}
