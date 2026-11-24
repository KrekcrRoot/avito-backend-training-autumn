package service

import (
	"avito-kitchen-service/internal/infrastructure/db/repository"
	"avito-kitchen-service/internal/infrastructure/db/sqlc"
	"avito-kitchen-service/internal/transport/http/dto"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) Create(ctx context.Context, request dto.RegisterUserRequest) (sqlc.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return sqlc.User{}, err
	}

	return s.userRepository.CreateUser(ctx, sqlc.CreateUserParams{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hash),
	})
}

func (s *UserService) GetByID(ctx context.Context, id int64) (sqlc.User, error) {
	return s.userRepository.GetUserByID(ctx, id)
}

func (s *UserService) Login(ctx context.Context, email, password string) (sqlc.User, error) {

	user, err := s.userRepository.GetUserByEmail(ctx, email)

	if err != nil {
		return sqlc.User{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil

}
