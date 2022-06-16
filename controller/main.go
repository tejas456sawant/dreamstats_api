package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tejas456sawant/dreamstats_api/config"
	"go.mongodb.org/mongo-driver/mongo"
)

var PlayerCollection *mongo.Collection = config.GetCollection(config.DB, "player")
var AllMatchesCollection *mongo.Collection = config.GetCollection(config.DB, "all")
var TopPlayersCollection *mongo.Collection = config.GetCollection(config.DB, "topPlayers")

func HelloWorld() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World!"})
	}
}
