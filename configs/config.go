package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	MConf MongoConfiguration
}

type MongoConfiguration struct {
	Username       string
	Password       string
	DbName         string
	WordCollection string
	UserCollection string
}

func GetConfig() Configuration {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file!")
	}

	// -------- Get env and assign
	get_username := os.Getenv("DB_USERNAME")
	get_password := os.Getenv("DB_PASSWORD")
	get_dbName := os.Getenv("DB_NAME")
	get_word_collection := os.Getenv("DB_COLLECTION_WORD")
	get_user_collection := os.Getenv("DB_COLLECTION_USER")

	mconf := &MongoConfiguration{
		Username:       get_username,
		Password:       get_password,
		DbName:         get_dbName,
		WordCollection: get_word_collection,
		UserCollection: get_user_collection,
	}
	configuration := Configuration{
		MConf: *mconf,
	}

	return configuration
}
