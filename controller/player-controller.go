package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

		id := c.Param("id")

		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		cache := rdb.Get(ctx, "player_"+id)

		if cache.Err() == redis.Nil {
			fmt.Println("cache miss")

			client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

			if err != nil {
				log.Fatal(err)
			}

			collection := client.Database("dreamstats").Collection("player")
			doc := collection.FindOne(context.Background(), bson.M{"cricinfo_id": id})

			var result bson.M
			doc.Decode(&result)

			var json_result []byte
			json_result, _ = json.Marshal(result)

			err = rdb.Set(ctx, "player_"+id, json_result, 24*time.Hour).Err()

			if err != nil {
				log.Fatal(err)
			}

			if result == nil {
				c.JSON(http.StatusNotFound, gin.H{"message": "Player not found."})
			} else {
				c.JSON(http.StatusOK, result)
			}
		} else {
			fmt.Println("cache hit")
			var result bson.M
			err := json.Unmarshal([]byte(cache.Val()), &result)
			if err != nil {
				log.Fatal(err)
			}
			c.JSON(http.StatusOK, result)
		}
	}
}
