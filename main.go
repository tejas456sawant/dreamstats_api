package main

import (
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/tejas456sawant/dreamstats_api/config"
	"github.com/tejas456sawant/dreamstats_api/controller"
)

func main() {
	router := gin.Default()

	config.ConnectDB()

	gin.SetMode(gin.ReleaseMode)

	router.SetTrustedProxies([]string{"localhost"})

	router.Static("/images", "/home/authorof_net/images")

	rdb := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "default",
		Password: "",
		DB:       0,
	}))

	router.GET("/", controller.HelloWorld())

	api_v1 := router.Group("/api/v1")
	{
		player := api_v1.Group("/player")
		{
			player.GET("/", controller.GetPlayer())
			player.GET("/:id", cache.CacheByRequestURI(rdb, 24*time.Hour), controller.GetPlayerById())
		}
		head := api_v1.Group("/head")
		{
			head.GET("/", cache.CacheByRequestURI(rdb, 24*time.Hour), controller.GetHeadToHead())
		}
		form := api_v1.Group("/form")
		{
			form.GET("/batting", cache.CacheByRequestURI(rdb, 24*time.Hour), controller.GetBattingForm())
			form.GET("/bowling", cache.CacheByRequestURI(rdb, 24*time.Hour), controller.GetBowlingForm())
		}
	}

	router.Run("localhost:8080")
}
