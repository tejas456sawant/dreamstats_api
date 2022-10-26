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
		DBURL = "mongodb+srv://doadmin:24E57t3LA19mnq0a@db-mongodb-blr1-93233-03fbf148.mongo.ondigitalocean.com"
	} else {
		DBURL = "mongodb+srv://doadmin:24E57t3LA19mnq0a@db-mongodb-blr1-93233-03fbf148.mongo.ondigitalocean.com"
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
