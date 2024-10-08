package config

import "os"

func GetMySQLDSN() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	return user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName
}
