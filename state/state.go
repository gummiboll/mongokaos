package state

import (
	"context"
	"log"
	"sync"

	"github.com/gummiboll/mongokaos/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	appState *types.AppState
	once     sync.Once
)

func initMongo(DBUrl string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DBUrl))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetAppState() *types.AppState {
	once.Do(func() {
		cfg := types.NewConfig()

		mongoClient, err := initMongo(cfg.DBUrl)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		log.Printf("Connected to MongoDB at %s\n", cfg.DBUrl)
		appState = &types.AppState{
			Config: cfg,
			Mongo:  mongoClient,
		}
	})
	return appState
}
