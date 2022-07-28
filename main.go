package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tejas456sawant/dreamstats_api/config"
	"github.com/tejas456sawant/dreamstats_api/controller"
)

func main() {
	router := gin.Default()

	config.ConnectDB()

	gin.SetMode(gin.ReleaseMode)

	router.SetTrustedProxies([]string{"localhost"})

	router.Static("/images", "/home/authorof_net/images")

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
		score := api_v1.Group("/score")
		{
			score.GET("/card/:id", controller.GetScorecardById())
			score.GET("/live", func(ctx *gin.Context) {
				ws, err := controller.Upgrade(ctx.Writer, ctx.Request)
				if err != nil {
					ctx.AbortWithError(http.StatusInternalServerError, err)
					return
				}
				controller.Writer(ws)
			})
		}
	}

	router.Run("localhost:8080")
}
