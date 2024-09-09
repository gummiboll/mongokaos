package state

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	ListenPort string
	APIKey     string
	DBUrl      string
	Debug      bool
}

func NewConfig() *Config {
	return &Config{
		ListenPort: os.Getenv("LISTEN_PORT"),
		APIKey:     os.Getenv("API_KEY"),
		DBUrl:      os.Getenv("DB_URL"),
		Debug:      os.Getenv("DEBUG") == "true",
	}
}

type AppState struct {
	Config *Config
	Mongo  *mongo.Client
}
