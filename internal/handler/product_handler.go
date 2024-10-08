package handler

import (
	"go-product/internal/adapter/dto"
	"go-product/internal/core/port"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service port.ProductService
}

func NewProductHandler(service port.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetProductById(c *gin.Context) {
	id := c.Param("id")
	product, err := h.service.GetById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, _ := h.service.GetAll()
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product dto.CreateProductRequest

	errorChecks := []func() error{
		func() error { return c.ShouldBindJSON(&product) },
		func() error { return product.Validate() },
		func() error { return h.service.Create(product) },
	}

	for _, check := range errorChecks {
		if err := check(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully created a new product"})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var product dto.UpdateProductRequest
	errorChecks := []func() error{
		func() error { return c.ShouldBindJSON(&product) },
		func() error { return product.Validate() },
		func() error { return h.service.Edit(c.Param("id"), product) },
	}

	for _, check := range errorChecks {
		if err := check(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated product"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
