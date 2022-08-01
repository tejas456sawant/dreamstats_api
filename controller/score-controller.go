package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetScorecardById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		id, err := strconv.ParseFloat(c.Param("id"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		doc := ScorecardCollection.FindOne(ctx, bson.M{"match_id": id})
		var result bson.M
		doc.Decode(&result)

		if result == nil {
			var url string
			url_check, _ := http.Get("https://www.espncricinfo.com/matches/engine/match/" + c.Param("id") + ".html")
			if strings.Contains(url_check.Request.URL.String(), "live-cricket-score") {
				url = strings.Replace(url_check.Request.URL.String(), "live-cricket-score", "full-scorecard", -1)
			} else {
				url = url_check.Request.URL.String()
			}
			defer url_check.Body.Close()

			res, err := http.Get(url)
			if err != nil {
				return
			}
			if res.StatusCode != 200 {
				return
			}
			defer res.Body.Close()

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Print(err)
			}

			script := doc.Find("script#__NEXT_DATA__").Text()

			var match map[string]map[string]map[string]map[string]interface{}
			json.Unmarshal([]byte(script), &match)

			var scorecard map[string]map[string]map[string]map[string]map[string]map[string][]interface{}
			json.Unmarshal([]byte(script), &scorecard)

			var playersOfTheMatch map[string]map[string]map[string]map[string]map[string]map[string]interface{}
			json.Unmarshal([]byte(script), &playersOfTheMatch)

			content := map[string]interface{}{
				"match_id":          id,
				"match":             match["props"]["appPageProps"]["data"]["match"],
				"scorecard":         scorecard["props"]["appPageProps"]["data"]["content"]["scorecard"],
				"playersOfTheMatch": playersOfTheMatch["props"]["appPageProps"]["data"]["content"]["supportInfo"]["playersOfTheMatch"],
			}

			if !strings.Contains(url_check.Request.URL.String(), "live-cricket-score") {
				ScorecardCollection.InsertOne(ctx, content)
			}

			c.JSON(http.StatusOK, content)
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}
