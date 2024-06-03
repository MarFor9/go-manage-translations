package api

import "go-template/internal/core/domain"

func toTranslationResponse(translation []domain.Translation) map[string]string {
	response := make(map[string]string)
	for _, t := range translation {
		response[t.Key] = t.Translation
	}
	return response
}
