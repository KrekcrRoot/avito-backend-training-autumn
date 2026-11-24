package service

import (
	"context"
	"restaurant-service/internal/infrastructure/db/repository"
	"restaurant-service/internal/infrastructure/db/sqlc"
)

type ProductService struct {
	repository *repository.ProductRepository
}

func NewProductService(repository *repository.ProductRepository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) GetProducts(ctx context.Context) ([]sqlc.Product, error) {
	return s.repository.GetAll(ctx)
}
