package config

import (
	"context"
	"log"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBURL string

func ConnectDB() *mongo.Client {
	if runtime.GOOS == "windows" {
		DBURL = "mongodb://root:dreamstats@dream.magiccup.store:27017"
	} else {
		DBURL = "mongodb://root:dreamstats@localhost:27017"
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(DBURL))
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
