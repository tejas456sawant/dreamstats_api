package main

import (
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/tejas456sawant/dreamstats_api/config"
	"github.com/tejas456sawant/dreamstats_api/controller"
)

func main() {
	router := gin.Default()

	config.ConnectDB()

	gin.SetMode(gin.ReleaseMode)

	router.SetTrustedProxies([]string{"localhost"})

	if runtime.GOOS == "windows" {
		router.Static("/images", "./images")
	} else {
		router.Static("/images", "/root/images")
	}

	// rdb := persist.NewRedisStore(redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Username: "default",
	// 	Password: "",
	// 	DB:       0,
	// }))

	// cache.CacheByRequestURI(rdb, 24*time.Hour),

	router.GET("/", controller.HelloWorld())

	api_v1 := router.Group("/api/v1")
	{
		player := api_v1.Group("/player")
		{
			player.GET("/", controller.GetPlayer())
			player.GET("/top", controller.GetTopPlayers())
			player.GET("/:id", controller.GetPlayerById())
		}
		head := api_v1.Group("/head")
		{
			head.GET("/", controller.GetHeadToHead())
		}
		form := api_v1.Group("/form")
		{
			form.GET("/batting", controller.GetBattingForm())
			form.GET("/bowling", controller.GetBowlingForm())
		}
		match := api_v1.Group("/match")
		{
			match.GET("/:id", controller.GetMatchById())
		}
		score := api_v1.Group("/scorecard")
		{
			score.GET("/:id", controller.GetScorecardById())
		}
	}

	router.Run("localhost:8080")
}
