package main

import (
	"context"
	"database/sql"
	"go-product/internal/adapter/db/mongodb"
	"go-product/internal/adapter/db/mysql"
	"go-product/internal/config"
	"go-product/internal/core/middleware"
	"go-product/internal/handler"
	"go-product/internal/routes"
	"go-product/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
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

	productRepository := mysql.NewProductRepositoryImpl(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	speedTestRepo := mongodb.NewSpeedTestRepository(mongoClient)
	speedTestService := service.NewSpeedTestService(speedTestRepo)

	router := gin.Default()
	router.Use(middleware.SpeedTestMiddleware(speedTestService))
	routes.RegisterProductRoutes(router, productHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
