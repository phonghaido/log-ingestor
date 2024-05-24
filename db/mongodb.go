package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/phonghaido/log-ingestor/data"
	"go.mongodb.org/mongo-driver/bson"
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

	document := bson.M{
		"level":      logData.Level,
		"message":    logData.Message,
		"resourceId": logData.ResourceID,
		"timestamp":  logData.Timestamp,
		"traceId":    logData.TraceID,
		"spanID":     logData.SpanID,
		"commit":     logData.Commit,
		"metadata": bson.M{
			"parentResourceId": logData.Metadata.ParentResourceID,
		},
	}
	_, err := coll.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

func Search(c echo.Context, client *mongo.Client, searchData data.SearchData) ([]data.LogData, error) {
	filter := bson.M{}

	log.Printf("Search data: %v\n", searchData)
	if searchData.Level != "" {
		filter["level"] = searchData.Level
	}
	if searchData.ResourceID != "" {
		filter["resourceId"] = searchData.ResourceID
	}
	if searchData.TraceID != "" {
		filter["traceId"] = searchData.TraceID
	}
	if searchData.SpanID != "" {
		filter["spanId"] = searchData.SpanID
	}
	var start time.Time
	var end time.Time
	var err error

	if searchData.StartDate != "" {
		start, err = time.Parse(time.RFC3339, searchData.StartDate)
		if err != nil {
			return make([]data.LogData, 0), err
		}
	} else {
		start = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	if searchData.EndDate != "" {
		end, err = time.Parse(time.RFC3339, searchData.EndDate)
		if err != nil {
			return make([]data.LogData, 0), err
		}
	} else {
		end = time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	filter["timestamp"] = bson.M{

		"$gte": start.Format(time.RFC3339),
		"$lte": end.Format(time.RFC3339),
	}

	coll := client.Database(os.Getenv("MONGODB_DATABASE")).Collection(os.Getenv("MONGODB_COLLECTION"))
	defer client.Disconnect(c.Request().Context())

	cursor, err := coll.Find(c.Request().Context(), filter)
	if err != nil {
		return make([]data.LogData, 0), err
	}
	defer cursor.Close(c.Request().Context())

	var results []data.LogData
	if err = cursor.All(c.Request().Context(), &results); err != nil {
		return make([]data.LogData, 0), err
	}
	return results, nil
}
