package main

// 中间件

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("example", "21232")
		// 让原本该执行的逻辑继续执行
		c.Next()

		end := time.Since(t)
		fmt.Printf("耗时:%V\n", end)
		status := c.Writer.Status()
		fmt.Println("状态:", status)
	}
}

func TokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		for k, v := range c.Request.Header {
			if k == "X-T oken" {
				token = v[0]
			}
		}
		if token != "bobby" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "未登录",
			})
			// 这里的return是起不到作用的,必须使用abort(),组织后续逻辑的执行
			// 为什么连return都阻止不了后续逻辑的执行
			// return
			c.Abort()
		}
		// 已登录
		c.Next()

	}
}
func main() {
	router := gin.New()
	// 在分组使用中间件
	//	authrized := router.Group("/goods")
	//传入，MyLogger()执行的结果
	router.Use(MyLogger())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run(":8083")

}
