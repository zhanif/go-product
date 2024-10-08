package routes

import (
	"go-product/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(router *gin.Engine, productHandler *handler.ProductHandler) {
	router.POST("/products", productHandler.CreateProduct)
	router.GET("/products/:id", productHandler.GetProductById)
	router.GET("/products", productHandler.GetAllProducts)
	router.PUT("/products/:id", productHandler.UpdateProduct)
	router.DELETE("/products/:id", productHandler.DeleteProduct)
}
