package handler

import (
	"avito-kitchen-service/internal/service"
	"avito-kitchen-service/internal/service/auth"
	"avito-kitchen-service/internal/transport/http/dto"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthHandler struct {
	JWTService  *auth.JWTService
	UserService *service.UserService
}

func NewAuthHandler(jwtService *auth.JWTService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		JWTService:  jwtService,
		UserService: userService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req dto.RegisterUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		WriteError(w, "INVALID_REQUEST", "invalid request data", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Create(r.Context(), req)

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			WriteError(w, "INVALID_REQUEST", "invalid request data", http.StatusBadRequest)
			return
		}

		WriteError(w, "INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	token, err := h.JWTService.GenerateToken(user.ID)

	if err != nil {
		WriteError(w, "INTERNAL_ERROR", "Failed to auth", http.StatusInternalServerError)
	}

	res := dto.RegisterUserResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	}

	WriteJSON(w, http.StatusOK, res)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, "INVALID_REQUEST", "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Login(context.Background(), req.Email, req.Password)

	if err != nil {
		WriteError(w, "INVALID_REQUEST", "invalid request", http.StatusUnauthorized)
		return
	}

	token, err := h.JWTService.GenerateToken(user.ID)

	if err != nil {
		WriteError(w, "INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
		return
	}

	res := dto.TokenResponse{
		Token: token,
	}

	WriteJSON(w, http.StatusOK, res)
}
