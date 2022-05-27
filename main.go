package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tejas456sawant/dreamstats_api/controller"
)

func main() {
	router := gin.Default()

	router.SetTrustedProxies([]string{"localhost"})

	router.Static("/images", "./images")

	router.GET("/", controller.HelloWorld())

	api_v1 := router.Group("/api/v1")
	{
		player := api_v1.Group("/player")
		{
			player.GET("/", controller.GetPlayer())
			player.GET("/:id", controller.GetPlayerById())
		}
		head := api_v1.Group("/head")
		{
			head.GET("/", controller.GetHeadToHead())
		}
	}

	router.Run("localhost:8080")
}
