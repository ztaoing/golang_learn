package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	protouser "golang_learn/golang_learn/gin_start/ch05/proto"
)

// 输出json
func main() {
	r := gin.Default()

	r.GET("/someJson", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hey",
			"status":  http.StatusOK,
		})
	})

	r.GET("moreJson", func(c *gin.Context) {
		// 也可以使用struct的输出json
		var msg struct {
			Name    string `json:"user"` // 此处Name会变为user
			Message string
			Number  int
		}
		msg.Name = "lena"
		msg.Message = "hey"
		msg.Number = 123
		c.JSON(http.StatusOK, &msg)

	})

	r.GET("/someProto", returnProto)
	r.Run(":8083")
}

func returnProto(c *gin.Context) {
	course := []string{
		"golang", "python",
	}
	user := &protouser.Teacher{
		Name:   "taking",
		Course: course,
	}
	c.ProtoBuf(http.StatusOK, user)
	// 这里返回的是proto串，如何反转义会可显示的信息？查看，proto章节
}
