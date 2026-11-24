package dto

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UserResponse struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}
