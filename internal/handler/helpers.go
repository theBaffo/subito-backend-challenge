package handler

import (
	"errors"
	"net/http"

	"github.com/theBaffo/subito-backend-challenge/internal/domain"
)

// errorResponse is the standard error envelope returned by all endpoints.
func errorResponse(message string) map[string]string {
	return map[string]string{"error": message}
}

// domainErrorToStatus maps known domain errors to their HTTP status codes.
// This centralises the translation in one place so handlers stay clean.
func domainErrorToStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrProductNotFound),
		errors.Is(err, domain.ErrOrderNotFound):
		return http.StatusNotFound

	case errors.Is(err, domain.ErrInvalidQuantity),
		errors.Is(err, domain.ErrEmptyOrder):
		return http.StatusUnprocessableEntity

	default:
		return http.StatusInternalServerError
	}
}
