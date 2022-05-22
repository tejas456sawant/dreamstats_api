package controller

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tejas456sawant/dreamstats_api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetHeadToHead() gin.HandlerFunc {
	return func(c *gin.Context) {
		batter := c.Query("batter")
		bowler := c.Query("bowler")

		ctx := context.Background()

		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := client.Disconnect(ctx); err != nil {
				log.Fatal(err)
			}
		}()

		coll := client.Database("dreamstats").Collection("ipl")

		output, _ := coll.Aggregate(ctx, bson.A{
			bson.D{{Key: "$match", Value: bson.D{
				{Key: "innings", Value: bson.D{
					{Key: "$elemMatch", Value: bson.D{
						{Key: "overs", Value: bson.D{
							{Key: "$elemMatch", Value: bson.D{
								{Key: "deliveries", Value: bson.D{
									{Key: "$elemMatch", Value: bson.D{
										{Key: "batter", Value: batter},
										{Key: "bowler", Value: bowler},
									}},
								}},
							}},
						}},
					}},
				}},
			}},
			},
			bson.D{
				{Key: "$addFields", Value: bson.D{
					{Key: "innings", Value: bson.D{{Key: "$map", Value: bson.D{
						{Key: "input", Value: "$innings"},
						{Key: "as", Value: "i"},
						{Key: "in", Value: bson.D{{Key: "overs", Value: bson.D{
							{Key: "$map", Value: bson.D{
								{Key: "input", Value: "$$i.overs"},
								{Key: "as", Value: "o"},
								{Key: "in", Value: bson.D{{Key: "deliveries", Value: bson.D{
									{Key: "$filter", Value: bson.D{
										{Key: "input", Value: "$$o.deliveries"},
										{Key: "as", Value: "b"},
										{Key: "cond", Value: bson.D{
											{Key: "$and", Value: bson.A{
												bson.D{{Key: "$eq", Value: bson.A{"$$b.batter", batter}}},
												bson.D{{Key: "$eq", Value: bson.A{"$$b.bowler", bowler}}},
												bson.D{{Key: "$lte", Value: bson.A{"$$b.extras.wides", false}}},
											}},
										}},
									}},
								}}}},
							}},
						}}}},
					}}}}}},
			},
			bson.D{{Key: "$sort", Value: bson.D{{Key: "info.dates", Value: 1}}}},
		})

		var results []bson.M
		if err = output.All(context.TODO(), &results); err != nil {
			panic(err)
		}

		var runs []int
		var balls []int
		var wickets []int
		var dots []int
		var sr []float64

		for _, result := range results {
			runs_local := 0
			wides_local := 0
			balls_local := 0
			wickets_local := 0
			dots_local := 0

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
			sr = append(sr, (float64(runs_local)/float64(balls_local))*100)
		}

		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"runs":    runs,
			"balls":   balls,
			"wickets": wickets,
			"dots":    dots,
			"sr":      sr,
			"total": gin.H{
				"runs":    utils.ArraySum(runs),
				"balls":   utils.ArraySum(balls),
				"wickets": utils.ArraySum(wickets),
				"dots":    utils.ArraySum(dots),
				"sr":      float64(utils.ArraySum(runs)) / float64(utils.ArraySum(balls)) * 100,
			},
		})
	}
}
