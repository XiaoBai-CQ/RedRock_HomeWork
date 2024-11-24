package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	engine.GET("/echo", func(c *gin.Context) {
		body, err := c.GetQuery("message")
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, gin.H{
			"message": body,
		})
	})
	engine.Run()
}
