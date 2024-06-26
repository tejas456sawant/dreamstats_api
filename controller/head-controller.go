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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetHeadToHead() gin.HandlerFunc {
	return func(c *gin.Context) {
		batter_id := c.Query("batter")
		bowler_id := c.Query("bowler")
		match_type := c.Query("match_type")
		group_by := c.Query("group_by")

		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		batter := PlayerCollection.FindOne(ctx, bson.M{"cricinfo_id": batter_id}, options.FindOne().SetProjection(bson.M{"identifier_name": 1}))
		bowler := PlayerCollection.FindOne(ctx, bson.M{"cricinfo_id": bowler_id}, options.FindOne().SetProjection(bson.M{"identifier_name": 1}))

		if batter.Err() != nil || bowler.Err() != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Head to head not found."})
			return
		}

		var batter_doc bson.M
		var bowler_doc bson.M
		batter.Decode(&batter_doc)
		bowler.Decode(&bowler_doc)
		batter_name := batter_doc["identifier_name"].(string)
		bowler_name := bowler_doc["identifier_name"].(string)

		query := queries.GetHeadQuery(batter_name, bowler_name, match_type, group_by)
		output, _ := AllMatchesCollection.Aggregate(ctx, query)

		var results []bson.M
		if err := output.All(ctx, &results); err != nil {
			panic(err)
		}

		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Head to head not found."})
		} else {
			var response []bson.M

			var global_runs int
			var global_balls int
			var global_wickets int
			var global_dots int

			for _, result := range results {
				var runs []int
				var balls []int
				var wickets []int
				var dots []int
				var sr []float64
				var ids []string
				var inns []int

				for _, matches := range result["matches"].(bson.A) {
					runs_local := 0
					wides_local := 0
					balls_local := 0
					wickets_local := 0
					dots_local := 0
					inns_local := 0

					for _, over := range matches.(bson.M)["innings"].(bson.A) {
						inns_local += 1
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

					global_balls += balls_local
					global_runs += runs_local
					global_wickets += wickets_local
					global_dots += dots_local

					dots = append(dots, dots_local)
					inns = append(inns, inns_local)
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
					"inns":       inns,
					"total": bson.M{
						"runs":    utils.ArraySum(runs),
						"balls":   utils.ArraySum(balls),
						"wickets": utils.ArraySum(wickets),
						"dots":    utils.ArraySum(dots),
						"inns":    utils.ArraySum(inns),
						"sr":      float64(utils.ArraySum(runs)) / float64(utils.ArraySum(balls)) * 100,
					},
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"data": response,
				"total": gin.H{
					"runs":    global_runs,
					"balls":   global_balls,
					"wickets": global_wickets,
					"dots":    global_dots,
				},
			})
		}
	}
}
