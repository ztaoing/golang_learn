package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func pong(c *gin.Context) {
	/*c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})*/
	// 等同于gin.H
	c.JSON(http.StatusOK, map[string]string{
		"message": "pong",
	})
}

func main() {
	r := gin.Default()
	r.GET("/ping", pong)
	r.Run(":8083") // listen and serve on 0.0.0.0:8080
}
