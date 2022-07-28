package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}

	return ws, nil
}

func Writer(ws *websocket.Conn) {
	for {
		ticker := time.NewTicker(time.Second * 3)
		defer ticker.Stop()
		for range ticker.C {
			res, err := http.Get("https://www.espncricinfo.com/live-cricket-score")
			if err != nil {
				log.Print(err)
				return
			}
			if res.StatusCode != 200 {
				log.Printf("status code error: %d %s", res.StatusCode, res.Status)
				return
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			var content map[string]map[string]map[string]map[string]map[string][]interface{}
			data := doc.Find("script#__NEXT_DATA__").Text()
			json.Unmarshal([]byte(data), &content)
			fmt.Println(content["props"]["appPageProps"]["data"]["content"]["matches"])

			err = ws.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func GetScorecardById() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.ParseFloat(c.Param("id"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		doc := ScorecardCollection.FindOne(context.Background(), bson.M{"match_id": id})
		var result bson.M
		doc.Decode(&result)

		if result == nil {
			res, err := http.Get("https://dreamstats-live-score.herokuapp.com/scorecard/" + c.Param("id"))
			if err != nil {
				fmt.Print(err)
				return
			}

			if res.StatusCode != 200 {
				fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
				return
			}

			var content map[string]interface{}
			json.NewDecoder(res.Body).Decode(&content)

			ScorecardCollection.InsertOne(context.Background(), content)

			c.JSON(http.StatusOK, content)
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}
