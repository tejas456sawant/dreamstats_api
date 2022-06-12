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

		// opts := options.Find().SetProjection(bson.D{{Key: "identifier_name", Value: 1}, {Key: "cricinfo_id", Value: 1}, {Key: "player.longName", Value: 1}, {Key: "_id", Value: 0}})
		docs, _ := PlayerCollection.Find(ctx, query)

		var result []bson.M
		docs.All(ctx, &result)

		// for i := range result {
		// 	result[i]["longName"] = result[i]["player"].(bson.M)["longName"]
		// 	delete(result[i], "player")
		// }

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
