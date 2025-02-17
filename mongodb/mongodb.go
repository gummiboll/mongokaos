package mongodb

import (
	"context"
	"errors"

	"github.com/gummiboll/mongokaos/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setFindOpts(reqData types.RequestData) *options.FindOptions {
	opts := options.Find()
	if reqData.Sort != nil {
		opts.SetSort(reqData.Sort)
	}
	if reqData.Limit != 0 {
		opts.SetLimit(reqData.Limit)
	}
	if reqData.Skip != 0 {
		opts.SetSkip(reqData.Skip)
	}
	if reqData.Projection != nil {
		opts.SetProjection(reqData.Projection)
	}
	return opts
}

func setFindOneOpts(reqData types.RequestData) *options.FindOneOptions {
	opts := options.FindOne()
	if reqData.Sort != nil {
		opts.SetSort(reqData.Sort)
	}
	if reqData.Projection != nil {
		opts.SetProjection(reqData.Projection)
	}
	return opts
}

// Define a type for MongoDB operations with flexible parameters
type MongoOperation func(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error)

// Function definitions
func findOne(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	var result bson.M
	err := collection.FindOne(ctx, reqData.Filter, setFindOneOpts(reqData)).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func findMany(ctx context.Context, collection *mongo.Collection, reqData types.RequestData) (interface{}, error) {
	return collection.Find(ctx, reqData.Filter, setFindOpts(reqData), nil)
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
