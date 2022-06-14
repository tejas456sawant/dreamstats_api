package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodb://root:dreamstats@api.fullbankdetails.com:27017
func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root1:dreamstats@api.fullbankdetails.com:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cancle()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

var DB *mongo.Client = ConnectDB()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("dreamstats").Collection(collectionName)
	return collection
}
