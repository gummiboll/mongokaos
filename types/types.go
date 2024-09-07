package types

import (
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type RequestData struct {
	DataSource string   `json:"dataSource"`
	Database   string   `json:"database"`
	Collection string   `json:"collection"`
	Filter     bson.D   `json:"filter"`
	Pipeline   []bson.D `json:"pipeline"`
	Document   bson.D   `json:"document"`
	Update     bson.D   `json:"update"`
	Upsert     bool     `json:"upsert"`
}

type SingleResult struct {
	Document primitive.M `json:"document,omitempty"`
}

type MultipleResults struct {
	Documents []primitive.M `json:"documents"`
}
