package repository

import (
	"avito-kitchen-service/internal/infrastructure/db/sqlc"
	"context"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(queries *sqlc.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]sqlc.User, error) {
	return r.queries.GetUsers(ctx)
}

func (r *UserRepository) CreateUser(ctx context.Context, args sqlc.CreateUserParams) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, args)
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (sqlc.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}
