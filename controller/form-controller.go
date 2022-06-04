package controller

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tejas456sawant/dreamstats_api/queries"
	"github.com/tejas456sawant/dreamstats_api/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func GetBattingForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		name := c.Query("name")
		limit, _ := strconv.ParseInt(c.Query("limit"), 0, 64)
		match_type := c.Query("match_type")

		query := queries.FormBatterQuery(name, limit, match_type)

		docs, _ := AllMatchesCollection.Aggregate(ctx, query)

		var results []bson.M
		docs.All(ctx, &results)

		var runs []int
		var balls []int
		var wickets []int
		var dots []int
		var sr []float64
		var my_team []string
		var opp_teams []string

		for _, result := range results {
			runs_local := 0
			wides_local := 0
			balls_local := 0
			wickets_local := 0
			dots_local := 0

			info := result["info"].(bson.M)
			players := info["players"].(bson.A)

			for _, player := range players {
				player_info := player.(bson.M)
				if !utils.CheckIfBsonContains(player_info["v"].(bson.A), name) {
					opp_teams = append(opp_teams, player_info["k"].(string))
				}

				if utils.CheckIfBsonContains(player_info["v"].(bson.A), name) {
					my_team = append(my_team, player_info["k"].(string))
				}
			}

			for _, over := range result["innings"].(bson.A) {
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
		}

		op := gin.H{
			"runs":      runs,
			"balls":     balls,
			"wickets":   wickets,
			"dots":      dots,
			"sr":        sr,
			"my_team":   my_team,
			"opp_teams": opp_teams,
		}

		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Form not found."})
		} else {
			c.JSON(http.StatusOK, op)
		}
	}
}

func GetBowlingForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		name := c.Query("name")
		limit, _ := strconv.ParseInt(c.Query("limit"), 0, 64)
		match_type := c.Query("match_type")

		query := queries.FormBowlerQuery(name, limit, match_type)

		docs, _ := AllMatchesCollection.Aggregate(ctx, query)

		var results []bson.M
		docs.All(ctx, &results)

		var runs []int
		var balls []int
		var wickets []int
		var dots []int
		var teams []string
		var opp_teams []string
		var economy []float64
		var average []float64

		for _, result := range results {
			runs_local := 0
			wides_local := 0
			balls_local := 0
			wickets_local := 0
			dots_local := 0

			info := result["info"].(bson.M)
			players := info["players"].(bson.A)

			for _, player := range players {
				player_info := player.(bson.M)
				if !utils.CheckIfBsonContains(player_info["v"].(bson.A), name) {
					opp_teams = append(opp_teams, player_info["k"].(string))
				}

				if utils.CheckIfBsonContains(player_info["v"].(bson.A), name) {
					teams = append(teams, player_info["k"].(string))
				}
			}

			for _, over := range result["innings"].(bson.A) {
				for _, delivery := range over.(bson.M) {
					for _, ball := range delivery.(bson.A) {
						for _, over := range ball.(bson.M) {
							for _, delivery := range over.(bson.A) {
								res := delivery.(bson.M)
								res2 := res["runs"].(bson.M)
								if res["extras"] != nil {
									extras := res["extras"].(bson.M)
									if extras["legbyes"] != nil || extras["byes"] != nil {
										runs_local += int(res2["batter"].(float64))
									} else {
										runs_local += int(res2["total"].(float64))
									}
								} else {
									runs_local += int(res2["total"].(float64))
								}

								if int(res2["total"].(float64)) == 0 {
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
			economy = append(economy, float64(runs_local)/(float64(balls_local)/6))

			if float64(wickets_local) == 0 {
				average = append(average, -1)
			} else {
				average = append(average, float64(runs_local)/float64(wickets_local))
			}
		}

		op := gin.H{
			"runs":      runs,
			"balls":     balls,
			"wickets":   wickets,
			"dots":      dots,
			"teams":     teams,
			"opp_teams": opp_teams,
			"economy":   economy,
			"average":   average,
			"total": gin.H{
				"runs":    utils.ArraySum(runs),
				"balls":   utils.ArraySum(balls),
				"wickets": utils.ArraySum(wickets),
				"dots":    utils.ArraySum(dots),
				"sr":      float64(float64(utils.ArraySum(runs))) / (float64(utils.ArraySum(balls) / 6)),
				"avg":     float64(utils.ArraySum(runs)) / float64(utils.ArraySum(wickets)),
			},
		}

		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Form not found."})
		} else {
			c.JSON(http.StatusOK, op)
		}
	}
}
