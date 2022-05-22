package controller

import "github.com/gin-gonic/gin"

func HelloWorld() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World!"})
	}
}
