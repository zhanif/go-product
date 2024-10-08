package port

import (
	"go-product/internal/adapter/dto"
	"go-product/internal/core/domain"
)

type ProductService interface {
	Create(product dto.CreateProductRequest) error
	GetById(id string) (*domain.Product, error)
	GetAll() ([]domain.Product, error)
	Edit(id string, product dto.UpdateProductRequest) error
	Delete(id string) error
}

type ProductRepository interface {
	FindById(id string) (*domain.Product, error)
	FindAll() ([]domain.Product, error)
	Save(product domain.Product) error
	Update(product domain.Product) error
	Destroy(id string) error
}
