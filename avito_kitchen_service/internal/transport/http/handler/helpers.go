package handler

import (
	"avito-kitchen-service/internal/transport/http/dto"
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

func WriteError(w http.ResponseWriter, code, message string, status int) {
	errRes := dto.ErrorResponse{
		Error: dto.Error{
			Code:    code,
			Message: message,
		},
	}
	WriteJSON(w, status, errRes)
}
