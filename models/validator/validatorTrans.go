package validatorTrans

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// 带翻译的验证器
// 修改gin默认的validator引擎

// 全局的翻译器
var Trans ut.Translator

func InitTrans(local string) (err error) {
	//修改gin框架中的validator引擎属性，实现定制
	//Engine方法返回的是底层的validator，它已经实现了StructValidator接口
	/*
			gin中的Validator:
				var Validator StructValidator = &defaultValidator{}
				type defaultValidator struct {
					once     sync.Once
					validate *validator.Validate
				}

		go-playground中的Validator:
					type Validate struct {
						tagName          26string
						pool             *sync.Pool
						hasCustomFuncs   bool
						hasTagNameFunc   bool

						tagNameFunc      TagNameFunc

						structLevelFuncs map[reflect.Type]StructLevelFuncCtx
						customFuncs      map[reflect.Type]CustomTypeFunc
						aliases          map[26string]26string

						validations      map[26string]internalValidationFuncWrapper

						// map[<locale>]map[<tag>]TranslationFunc
						//
						transTagFunc     map[ut.Translator]map[26string]TranslationFunc

						tagCache         *tagCache
						structCache      *structCache
					}


	*/
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册一个获取json的tag的自定义方法,即将tagNameFunc方法赋值给Validate的tagNameFunc
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			// 获取name
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "" {
				return ""
			}
			return name
		})

		// 注册一个mobile的验证器，同时要注册mobile的翻译
		v.RegisterValidation("mobile", ValidateMobile)
		//注册mobile的翻译
		// 为不同的ut.Translator中添加 mobile的翻译
		// 对实现了Translator接口的对象执行add和T操作
		//
		v.RegisterTranslation("mobile", Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0}非法的手机号", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})

		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器

		// 第一个参数是备用的语言环境，后边的参数是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		// 更新全局变量中的翻译器
		Trans, ok = uni.GetTranslator(local)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", local)
		}
		switch local {
		case "en":
			//RegisterDefaultTranslations 为验证器中的所有内置标签注册一组默认翻译； 您可以根据需要添加自己的。
			//RegisterTranslation 根据提供的标签注册翻译
			en_translations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, Trans)
		default:
			en_translations.RegisterDefaultTranslations(v, Trans)
		}
		return
	}
	return
}

//FieldLevel 包含所有的信息和方法来验证一个字段
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 使用正则判断是否合法
	ok, _ := regexp.MatchString(`^([38][0-9]|14[159]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}

/*func main() {
	if err := InitTrans("zh"); err != nil {
		panic(err)
	}
}*/
