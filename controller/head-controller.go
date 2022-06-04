package controller

import (
	"context"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tejas456sawant/dreamstats_api/queries"
	"github.com/tejas456sawant/dreamstats_api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetHeadToHead() gin.HandlerFunc {
	return func(c *gin.Context) {
		batter := c.Query("batter")
		bowler := c.Query("bowler")
		match_type := c.Query("match_type")
		group_by := c.Query("group_by")

		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		query := queries.GetHeadQuery(batter, bowler, match_type, group_by)
		output, _ := AllMatchesCollection.Aggregate(ctx, query)

		var results []bson.M
		if err := output.All(ctx, &results); err != nil {
			panic(err)
		}

		var response []bson.M

		for _, result := range results {
			var runs []int
			var balls []int
			var wickets []int
			var dots []int
			var sr []float64
			var ids []string

			for _, matches := range result["matches"].(bson.A) {
				runs_local := 0
				wides_local := 0
				balls_local := 0
				wickets_local := 0
				dots_local := 0

				for _, over := range matches.(bson.M)["innings"].(bson.A) {
					for _, delivery := range over.(bson.M) {
						for _, ball := range delivery.(bson.A) {
							for _, over := range ball.(bson.M) {
								for _, delivery := range over.(bson.A) {
									res := delivery.(bson.M)
									res2 := res["runs"].(bson.M)
									runs_local += int(res2["batter"].(float64))

									if int(res2["batter"].(float64)) == 0 {
										dots_local += 1
									}

									if res["extras"] != nil {
										extras := res["extras"].(bson.M)
										if extras["wides"] != nil {
											wides_local += 1
										}
										if extras["noballs"] != nil {
											wides_local += 1
										}
									}

									if res["wickets"] != nil {
										for _, wicket := range res["wickets"].(bson.A) {
											if wicket.(bson.M)["kind"] != "run out" || wicket.(bson.M)["kind"] != "retired hurt" {
												wickets_local += 1
											}
										}
									}

									balls_local += 1
								}
							}
						}
					}
				}
				balls_local = balls_local - wides_local

				runs = append(runs, runs_local)
				balls = append(balls, balls_local)
				wickets = append(wickets, wickets_local)
				dots = append(dots, dots_local)
				if math.IsNaN((float64(runs_local) / float64(balls_local)) * 100) {
					sr = append(sr, 0)
				} else {
					sr = append(sr, (float64(runs_local)/float64(balls_local))*100)
				}
				ids = append(ids, matches.(bson.M)["_id"].(primitive.ObjectID).Hex())
			}

			var group_by_local string = ""
			var match_type_local string = ""
			if result["_id"] != nil {
				id := result["_id"].(bson.M)

				if id["match_type"] != nil {
					match_type_local = id["match_type"].(string)
				}

				if id["group_by"] != nil {
					if reflect.ValueOf(id["group_by"]).Kind() == reflect.Float64 {
						group_by_local = strconv.FormatFloat(id["group_by"].(float64), 'f', -1, 64)
					}
					if reflect.ValueOf(id["group_by"]).Kind() == reflect.String {
						group_by_local = id["group_by"].(string)
					}
				}

			}

			response = append(response, bson.M{
				"group_by":   group_by_local,
				"match_type": match_type_local,
				"ids":        ids,
				"runs":       runs,
				"balls":      balls,
				"wickets":    wickets,
				"dots":       dots,
				"sr":         sr,
				"total": bson.M{
					"runs":    utils.ArraySum(runs),
					"balls":   utils.ArraySum(balls),
					"wickets": utils.ArraySum(wickets),
					"dots":    utils.ArraySum(dots),
					"sr":      float64(utils.ArraySum(runs)) / float64(utils.ArraySum(balls)) * 100,
				},
			})

		}

		c.JSON(http.StatusOK, response)
	}
}
