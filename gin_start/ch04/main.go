package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 声明约束
type Person struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/user/:name/:action", func(c *gin.Context) {
		var person Person
		// 绑定限制
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		name := c.Param("name")
		action := c.Param("action")
		c.String(http.StatusOK, "%s is %s", name, action)
	})

	r.GET("/usr/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")

		c.String(http.StatusOK, "%s is %s", name, action)
	})
	r.Run(":8082")
}
