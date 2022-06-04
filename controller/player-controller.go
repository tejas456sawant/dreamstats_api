package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

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

		docs, _ := PlayerCollection.Find(ctx, query)

		var result []bson.M
		docs.All(ctx, &result)

		if len(result) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Player not found."})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func GetPlayerById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		id := c.Param("id")

		doc := PlayerCollection.FindOne(ctx, bson.M{"cricinfo_id": id})

		var result bson.M
		doc.Decode(&result)

		if result == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Player not found."})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}
