package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tejas456sawant/dreamstats_api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMatchById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		if c.Param("id") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		doc := AllMatchesCollection.FindOne(ctx, bson.M{"_id": id})

		if doc.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": doc.Err().Error()})
			return
		}

		var result models.Match
		doc.Decode(&result)

		if result.ID.String() == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Match not found."})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}
