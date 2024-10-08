package mongodb

import (
	"context"
	"go-product/internal/config"
	"go-product/internal/core/port"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type SpeedTestRepositoryImpl struct {
	collection *mongo.Collection
}

func NewSpeedTestRepository(client *mongo.Client) port.SpeedTestRepository {
	collection := client.Database(config.GetMongoDBName()).Collection("speed_tests")
	return &SpeedTestRepositoryImpl{collection: collection}
}

func (r *SpeedTestRepositoryImpl) Save(method, path string, duration time.Duration) error {
	doc := map[string]interface{}{
		"method":    method,
		"path":      path,
		"duration":  duration.Seconds(),
		"timestamp": time.Now(),
	}

	_, err := r.collection.InsertOne(context.TODO(), doc)
	return err
}
