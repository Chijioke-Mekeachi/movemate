package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var ShipmentCollection *mongo.Collection

func init() {
	_ = ConnectMongo()
}

func ConnectMongo() error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "ticketdb"
	}

	collName := os.Getenv("MONGO_SHIPMENT_COLLECTION")
	if collName == "" {
		collName = "shipments"
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	MongoClient = client
	ShipmentCollection = client.Database(dbName).Collection(collName)
	return nil
}
