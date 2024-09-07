package mongodb

import (
	"context"
	"errors"

	"github.com/gummiboll/mongokaos/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define a type for MongoDB operations with flexible parameters
type MongoOperation func(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error)

// Function definitions
func findOne(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	return collection.FindOne(ctx, reqData.Filter), nil
}

func findMany(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	return collection.Find(ctx, reqData.Filter, options.Find())
}

func aggregate(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	return collection.Aggregate(ctx, reqData.Pipeline)
}

func updateOne(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	opts := options.Update().SetUpsert(reqData.Upsert)
	return collection.UpdateOne(ctx, reqData.Filter, reqData.Update, opts)
}

func updateMany(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	opts := options.Update().SetUpsert(reqData.Upsert)
	return collection.UpdateMany(ctx, reqData.Filter, reqData.Update, opts)
}

func insertOne(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	return collection.InsertOne(ctx, reqData.Document)
}

func deleteOne(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	return collection.DeleteOne(ctx, reqData.Filter)
}

func deleteMany(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	return collection.DeleteMany(ctx, reqData.Filter)
}

// Define the function map
var FunctionMap = map[string]MongoOperation{
	"findOne":    findOne,
	"find":       findMany,
	"aggregate":  aggregate,
	"updateOne":  updateOne,
	"updateMany": updateMany,
	"insertOne":  insertOne,
	"deleteOne":  deleteOne,
	"deleteMany": deleteMany,
}

// ExecuteAction is used to run the appropriate MongoDB operation
func ExecuteAction(action string, ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	mongoFunc, exists := FunctionMap[action]
	if !exists {
		return nil, errors.New("action not found")
	}
	return mongoFunc(ctx, collection, reqData)
}
