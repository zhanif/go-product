package config

import (
	"os"
)

func GetMongoDBURI() string {
	user := os.Getenv("MONGODB_USER")
	password := os.Getenv("MONGODB_PASSWORD")
	host := os.Getenv("MONGODB_HOST")
	port := os.Getenv("MONGODB_PORT")
	dbName := GetMongoDBName()

	return "mongodb://" + user + ":" + password + "@" + host + ":" + port + "/" + dbName
}

func GetMongoDBName() string {
	return os.Getenv("MONGODB_DB_NAME")
}
