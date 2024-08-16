package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	MongoURI    string
	MongoDBName string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("/home/saidazim/NT_homeworks/91_homework/.env")
	if err != nil {
		log.Printf(".env fayli yuklanmadi: %v", err)
	}

	viper.AutomaticEnv()

	return &Config{
		Port:        viper.GetString("PORT"),
		MongoURI:    viper.GetString("MONGO_URI"),
		MongoDBName: viper.GetString("MONGO_DB_NAME"),
	}, nil
}
