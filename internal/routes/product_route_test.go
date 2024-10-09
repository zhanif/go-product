package routes

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-product/internal/adapter/db/mongodb"
	"go-product/internal/adapter/db/mysql"
	"go-product/internal/adapter/dto"
	"go-product/internal/config"
	"go-product/internal/core/domain"
	"go-product/internal/core/middleware"
	"go-product/internal/handler"
	"go-product/internal/service"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Unable to load environment variables")
	}

	dsn := config.GetMySQLDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	mongoURI := config.GetMongoDBURI()
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	productRepository := mysql.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	speedTestRepo := mongodb.NewSpeedTestRepository(mongoClient)
	speedTestService := service.NewSpeedTestService(speedTestRepo)

	router.Use(middleware.SpeedTestMiddleware(speedTestService))
	RegisterProductRoutes(router, productHandler)

	return router
}

func getFirstProductID(router *gin.Engine, t *testing.T) string {
	req, _ := http.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var products []domain.Product
	if err := json.Unmarshal(w.Body.Bytes(), &products); err != nil || len(products) == 0 {
		t.Fatalf("Unable to fetch product data\n%v", err)
	}
	return products[0].ID
}

func performRequest(router *gin.Engine, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestGetAllProducts_Empty(t *testing.T) {
	router := setupRouter()
	w := performRequest(router, http.MethodGet, "/products", nil)

	var responseBody []domain.Product
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Empty(t, responseBody)
}

func TestGetAllProducts_NotEmpty(t *testing.T) {
	router := setupRouter()
	w := performRequest(router, http.MethodGet, "/products", nil)

	var responseBody []domain.Product
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseBody)
}

func TestGetProductById_Success(t *testing.T) {
	router := setupRouter()
	id := getFirstProductID(router, t)
	w := performRequest(router, http.MethodGet, fmt.Sprintf("/products/%s", id), nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProductById_Fail(t *testing.T) {
	router := setupRouter()
	id := "invalid-id"
	w := performRequest(router, http.MethodGet, fmt.Sprintf("/products/%s", id), nil)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestCreateProduct_Success(t *testing.T) {
	router := setupRouter()
	product := dto.CreateProductRequest{
		Name:  "Test Product",
		Price: 10000,
		Stock: 10,
	}

	jsonValue, _ := json.Marshal(product)
	w := performRequest(router, http.MethodPost, "/products", jsonValue)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateProduct_Fail(t *testing.T) {
	router := setupRouter()
	product := dto.CreateProductRequest{
		Name:  "Test Product",
		Price: 0,
	}

	jsonValue, _ := json.Marshal(product)
	w := performRequest(router, http.MethodPost, "/products", jsonValue)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestUpdateProduct_Success(t *testing.T) {
	router := setupRouter()
	id := getFirstProductID(router, t)

	product := dto.UpdateProductRequest{
		Name:  "Updated Test Product",
		Price: 10_000,
		Stock: 10,
	}

	jsonValue, _ := json.Marshal(product)
	w := performRequest(router, http.MethodPut, fmt.Sprintf("/products/%s", id), jsonValue)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateProduct_Fail(t *testing.T) {
	router := setupRouter()
	id := getFirstProductID(router, t)

	product := dto.UpdateProductRequest{
		Name: "-",
	}

	jsonValue, _ := json.Marshal(product)
	w := performRequest(router, http.MethodPut, fmt.Sprintf("/products/%s", id), jsonValue)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestDeleteProduct_Success(t *testing.T) {
	router := setupRouter()
	id := getFirstProductID(router, t)
	w := performRequest(router, http.MethodDelete, fmt.Sprintf("/products/%s", id), nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteProduct_Fail(t *testing.T) {
	router := setupRouter()
	id := "invalid-id"
	w := performRequest(router, http.MethodDelete, fmt.Sprintf("/products/%s", id), nil)
	assert.NotEqual(t, http.StatusOK, w.Code)
}
