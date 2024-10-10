package service

import (
	"go-product/internal/adapter/dto"
	"go-product/internal/core/domain"
	"go-product/internal/core/port"
	"log"

	"github.com/mitchellh/mapstructure"
)

type ProductServiceImpl struct {
	repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) port.ProductService {
	return &ProductServiceImpl{repo: repo}
}

func (s *ProductServiceImpl) Create(productDto dto.CreateProductRequest) error {
	var product domain.Product
	if err := mapstructure.Decode(productDto, &product); err != nil {
		log.Fatal("Error: ", err)
		return err
	}

	return s.repo.Save(product)
}

func (s *ProductServiceImpl) GetById(id string) (*domain.Product, error) {
	return s.repo.FindById(id)
}

func (s *ProductServiceImpl) GetAll() ([]domain.Product, error) {
	return s.repo.FindAll()
}

func (s *ProductServiceImpl) Edit(id string, productDto dto.UpdateProductRequest) error {
	var product domain.Product
	if err := mapstructure.Decode(productDto, &product); err != nil {
		log.Fatal("Error: ", err)
		return err
	}

	return s.repo.Update(id, product)
}

func (s *ProductServiceImpl) Delete(id string) error {
	return s.repo.Destroy(id)
}
