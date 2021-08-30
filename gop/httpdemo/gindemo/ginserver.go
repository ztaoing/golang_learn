package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const keyRequestId = "requestId"

func main() {
	r := gin.Default()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	r.Use(func(c *gin.Context) {
		// 用来记录日志的中间件

		// 记录任务的开始时间
		s := time.Now()
		// 在next中才真正的执行
		//继续执行,在每一次的请求都会执行这个logger
		c.Next()
		// 结束时间
		e := time.Now()
		logger.Info("incoming request",
			zap.String("path", c.Request.URL.Path), // 路径
			zap.Int("status", c.Writer.Status()),   //状态
			zap.Duration("lantency", e.Sub(s)),     // 处理所用的时间
		)
	}, func(c *gin.Context) {
		// 用来在每个请求中添加请求id的中间件
		c.Set(keyRequestId, rand.Int())
		c.Next()
	})
	r.GET("/ping", func(c *gin.Context) {
		h := gin.H{
			"message": "pong",
		}
		if id, exsits := c.Get(keyRequestId); exsits {
			// 获得到id的时候才添加到响应中
			h[keyRequestId] = id
		}

		c.JSON(200, h)
	})
	r.Run()
}
