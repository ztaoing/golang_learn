package main

import "github.com/gin-gonic/gin"

func main() {
	//  使用默认中间件创建一个gin路由器
	// Default默认开启 logger and recover 中间件
	// 	engine.Use(Logger(), Recovery())
	router := gin.Default()

	// router := gin.New() // 没有添加任何默认中间件

	// restful 的开发
	router.GET("/get", getting)
	router.POST("/post", posting)
	router.PUT("/put", puting)
	router.PATCH("/patch", patching)
	router.HEAD("/head", heading)
	router.OPTIONS("/options", options)

}

func options(context *gin.Context) {

}

func heading(context *gin.Context) {

}

func patching(context *gin.Context) {

}

func puting(context *gin.Context) {

}

func posting(context *gin.Context) {

}

func getting(context *gin.Context) {

}
