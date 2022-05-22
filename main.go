package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tejas456sawant/dreamstats_api/controller"
)

func main() {
	router := gin.Default()

	router.SetTrustedProxies([]string{"localhost"})

	router.GET("/", controller.HelloWorld())

	api_v1 := router.Group("/api/v1")
	{
		player := api_v1.Group("/player")
		{
			player.GET("/", controller.GetPlayer())
		}
	}

	router.Run("localhost:8080")
}
