package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type player struct {
	cricinfo_id     string
	identifier      string
	identifier_name string
	player          interface{}
	content         struct {
		teams          []interface{}
		profile        interface{}
		careerAverages interface{}
		topRecords     []interface{}
	}
}

func GetPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		query := bson.M{}
		gender := c.Query("gender")
		name := c.Query("name")

		if gender != "" {
			query = bson.M{"player.gender": gender}
		}
		if name != "" {
			query = bson.M{"$text": bson.M{"$search": name}}
		}
		if name != "" && gender != "" {
			query = bson.M{"$text": bson.M{"$search": name}, "player.gender": gender}
		}

		opts := options.Find().SetProjection(bson.D{{Key: "content", Value: 0}})
		docs, _ := PlayerCollection.Find(ctx, query, opts)

		var result []bson.M
		docs.All(ctx, &result)

		if len(result) == 0 {
			c.JSON(http.StatusOK, []bson.M{})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func GetPlayerById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		id := c.Param("id")

		doc := PlayerCollection.FindOne(ctx, bson.M{"cricinfo_id": id})

		var result bson.M
		doc.Decode(&result)

		if result == nil {
			res, err := http.Get("https://www.espncricinfo.com/ci/content/player/" + id + ".html")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
				return
			}
			if res.StatusCode != 200 {
				c.JSON(http.StatusNotFound, gin.H{"message": "Player not found."})
				return
			}
			defer res.Body.Close()

			pageDoc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Print(err)
			}

			script := pageDoc.Find("script#__NEXT_DATA__").Text()

			var player_data map[string]map[string]map[string]map[string]map[string]interface{}
			json.Unmarshal([]byte(script), &player_data)

			var content map[string]map[string]map[string]map[string]map[string][]interface{}
			var content2 map[string]map[string]map[string]map[string]map[string]interface{}

			json.Unmarshal([]byte(script), &content)
			json.Unmarshal([]byte(script), &content2)

			player := player{
				cricinfo_id:     id,
				identifier:      "",
				identifier_name: player_data["props"]["appPageProps"]["data"]["player"]["name"].(string),
			}

			player.player = player_data["props"]["appPageProps"]["data"]["player"]

			player.content.teams = content["props"]["appPageProps"]["data"]["content"]["teams"]
			player.content.profile = content2["props"]["appPageProps"]["data"]["content"]["profile"]
			player.content.careerAverages = content2["props"]["appPageProps"]["data"]["content"]["careerAverages"]
			player.content.topRecords = content["props"]["appPageProps"]["data"]["content"]["topRecords"]

			insertDoc := bson.D{
				{Key: "cricinfo_id", Value: player.cricinfo_id},
				{Key: "identifier", Value: player.identifier},
				{Key: "identifier_name", Value: player.identifier_name},
				{Key: "player", Value: player.player},
				{Key: "content", Value: bson.D{
					{Key: "teams", Value: player.content.teams},
					{Key: "profile", Value: player.content.profile},
					{Key: "careerAverages", Value: player.content.careerAverages},
					{Key: "topRecords", Value: player.content.topRecords},
				}},
			}

			PlayerCollection.InsertOne(ctx, insertDoc)

			var result bson.M
			doc = PlayerCollection.FindOne(ctx, bson.M{"cricinfo_id": id})
			doc.Decode(&result)

			c.JSON(http.StatusOK, result)
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func GetTopPlayers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()

		docs, _ := TopPlayersCollection.Find(ctx, bson.M{})

		var result []bson.M
		docs.All(ctx, &result)

		if len(result) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Player not found."})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}
