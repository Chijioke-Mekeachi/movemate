package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var ShipmentCollection *mongo.Collection

func init() {
	if err := ConnectMongo(); err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
}

func isPlaceholderURI(uri string) bool {
	lower := strings.ToLower(uri)
	return strings.Contains(lower, "<password>") || strings.Contains(lower, "<username>") || strings.Contains(lower, "<user>") || strings.Contains(lower, "<host>")
}

func ConnectMongo() error {
	uri := os.Getenv("MONGO_URI")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")

	if uri != "" && isPlaceholderURI(uri) {
		if mongoUser != "" && mongoPassword != "" && mongoHost != "" {
			uri = fmt.Sprintf("mongodb+srv://%s:%s@%s/?appName=Movemate", mongoUser, mongoPassword, mongoHost)
		} else {
			return fmt.Errorf("invalid MONGO_URI placeholder detected; set a valid MONGO_URI or provide MONGO_USER, MONGO_PASSWORD, and MONGO_HOST")
		}
	}

	if uri == "" {
		if mongoUser != "" && mongoPassword != "" && mongoHost != "" {
			uri = fmt.Sprintf("mongodb+srv://%s:%s@%s/?appName=Movemate", mongoUser, mongoPassword, mongoHost)
		} else {
			uri = "mongodb://localhost:27017"
		}
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
