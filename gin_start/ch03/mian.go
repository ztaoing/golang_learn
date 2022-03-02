package main

import "github.com/gin-gonic/gin"

// 路由分组: 例如 提供商品服务 又要提供用户服务
func main() {
	router := gin.Default()
	// 新建分组
	goodsGroup := router.Group("goods")
	{
		goodsGroup.GET("/05list", goodsList)
		goodsGroup.POST("/add", createGoods)
	}

	/*
		v1 := router.Group("v1")
		{

		}

		v2 := router.Group("v2")
		{

		}*/

	router.Run(":8082")
}

func goodsList(context *gin.Context) {

}

func createGoods(context *gin.Context) {

}

func goodsDetails(context *gin.Context) {

}
