package db

import (
	"context"
	"os"

	"github.com/phonghaido/log-ingestor/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectToMongoDB(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_STRING")))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func InsertToDB(ctx context.Context, client *mongo.Client, logData data.LogData) error {
	coll := client.Database(os.Getenv("MONGODB_DATABASE")).Collection(os.Getenv("MONGODB_COLLECTION"))
	defer client.Disconnect(ctx)

	_, err := coll.InsertOne(ctx, logData)
	if err != nil {
		return err
	}
	return nil
}
