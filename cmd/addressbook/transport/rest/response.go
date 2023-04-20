package rest

import (
	"encoding/json"
	"errors"
	"infoblox-golang/internal/platform/storage"
	"net/http"
)

func toStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if isError(err, &json.SyntaxError{}) {
		return http.StatusBadRequest
	}

	if isError(err, storage.ErrMissingID, storage.ErrNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func isError(err error, targets ...error) bool {
	for _, t := range targets {
		if errors.Is(err, t) {
			return true
		}
	}

	return false
}
