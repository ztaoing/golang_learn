package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	abs, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(abs)
	//模板的路径,这个目录是一个相对目录
	router.LoadHTMLGlob("templates/**/*")
	//router.LoadHTMLFiles(abs+"templates/index.tmpl", "templates/goods.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "中国",
		})
	})

	// 如果没有在模板中定义define，那么就可以使用默认的文件夹名来查找，可以解决文件名冲突的问题
	router.GET("/goods/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "goods/list.html", gin.H{
			"name": "微服务",
		})
	})

	router.GET("/user/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user/list.html", gin.H{
			"name": "微服务",
		})
	})
}
