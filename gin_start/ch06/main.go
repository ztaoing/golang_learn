package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/validator/v10"
)

// 表单验证

type Login struct {
	User     string `form:"user1" json:"user2" xml:"user3" binding:"required,min=3,max=10"` // form需要使用user1，json需要使用user2，xml需要使用user3
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type SignUpParam struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}

	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// 更换翻译器
func InitTrans(locale string) (err error) {
	// 修改gin框架中的validator引擎属性，实现定制
	// 把gin的validator变成go-playground中的validator.Validate
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 获取一个获取json的tag的自定义方法
		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "_" {
				return ""
			}
			return name
		})
		/// 变换成功之后就需要按照go-playground中validator来做了
		zhT := zh.New() //中文翻译器
		enT := en.New()
		// 第一个参数是备用的语言环境，后边的参数是应该支持的语言环境
		uni := ut.New(enT, zhT)

		// 修改全局的trans
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) not found", locale)
		}

		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(validate, trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(validate, trans)
		default:
			en_translations.RegisterDefaultTranslations(validate, trans)
		}
		return
	}
	return
}

// 定义一个全局的翻译器
var trans ut.Translator

func main() {
	if err := InitTrans("zh"); err != nil {
		fmt.Println("初始化翻译器错误")
		return
	}

	r := gin.Default()
	r.POST("/loginJSON", func(c *gin.Context) {
		var loginForm Login

		if err := c.ShouldBind(&loginForm); err != nil {
			// 进一步处理拿到的err
			errors, ok := err.(validator.ValidationErrors)
			if !ok {
				// 如果转换错误，就返回转换错误的信息，
				c.JSON(http.StatusOK, gin.H{
					"msg": errors.Error(),
				})
			}

			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"err": removeTopStruct(errors.Translate(trans)), // 使用定义的翻译器进行翻译
			})
			return
		}
		// 没有错误
		c.JSON(http.StatusOK, gin.H{
			"msg": "登录成功",
		})

	})

	r.POST("/signup", func(c *gin.Context) {
		var signupForm SignUpParam
		if err := c.ShouldBind(&signupForm); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
				"age": c.Param("age"),
			})
			return
		}
		// 没有错误
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
		})

	})

	_ = r.Run(":8083")
}
