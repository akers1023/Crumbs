package util

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func HandleError(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
