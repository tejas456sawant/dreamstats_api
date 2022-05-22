package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		if err != nil {
			log.Fatal(err)
		}

		query := bson.M{}
		gender := c.Query("gender")
		name := c.Query("name")

		if gender != "" {
			query = bson.M{"player.gender": gender}
		}
		if name != "" {
			query = bson.M{"$text": bson.M{"$search": name}}
		}
		if name != "" && gender != "" {
			query = bson.M{"$text": bson.M{"$search": name}, "player.gender": gender}
		}

		collection := client.Database("dreamstats").Collection("player")
		docs, _ := collection.Find(context.Background(), query)

		var result []bson.M
		docs.All(context.Background(), &result)

		if len(result) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Player not found."})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func GetPlayerById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		if err != nil {
			log.Fatal(err)
		}

		id := c.Param("id")

		collection := client.Database("dreamstats").Collection("player")
		doc := collection.FindOne(context.Background(), bson.M{"cricinfo_id": id})

		var result bson.M
		doc.Decode(&result)

		if result == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Player not found."})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}
