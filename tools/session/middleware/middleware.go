package middleware

import (
	"errors"
	"golang_learn/golang_learn/tools/config"
	"golang_learn/golang_learn/tools/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token is valid yet")
	TokenMalformed   = errors.New("token is  inllegal") // token不合法
	TokenInvalid     = errors.New("can not handle this token")
)

// 验证用户：从header中获得到x-token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// jwt的鉴权，通过获取header中的x-token。登录时，返回的token信息，
		//这里前端需要把token存储到cookie或者本地的存储中，
		//不过需要和后端协商过期时间，可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		// 没有token，即没有登录
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}
		// 有token，也不一定是正确的，所以需要解析后校验
		j := NewJWT()
		// 解析token中的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			// 如果是token过期
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "授权已过期",
				})
				c.Abort()
				return
			}
			//token的其他情况
			c.JSON(http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}

		// 解析成功
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next() // 执行之后的操作
	}
}

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(config.SigningKey), // 可以为jwt设置过期时间，使用过期的jwt是不可以的，这里只是简单的配置，所以没有加入过期时间
	}
}

/**
type Token struct {
	Raw       string                 // The raw token.  Populated when you Parse a token
	Method    SigningMethod          // The signing method used or to be used
	Header    map[string]interface{} // The first segment of the token
	Claims    Claims                 // The second segment of the token
	Signature string                 // The third segment of the token.  Populated when you Parse a token
	Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}
*/

// 解析token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	// 使用jwt组件解析jwt
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	// 解析出现错误
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {

			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

// 创建token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	//给TimeFunc设置默认值
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		//断言成功，覆盖默认值，设置为当前时间
		jwt.TimeFunc = time.Now
		//设置过期时间
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour * 1).Unix()
		//根据claims重新生成token
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

/**
type Context struct {
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter

	Params   Params
	handlers HandlersChain
	index    int8
	fullPath string

	engine *Engine
	params *Params

	// This mutex protect Keys map
	mu sync.RWMutex

	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]interface{}

	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
	Errors errorMsgs

	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string

	// queryCache use url.ParseQuery cached the param query result from c.Request.URL.Query()
	queryCache url.Values

	// formCache use url.ParseQuery cached PostForm contains the parsed form data from POST, PATCH,
	// or PUT body parameters.
	formCache url.Values

	// SameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests.
	sameSite http.SameSite
}
*/
/**
将index设置为末尾，即结束。
这里就涉及到index的作用！
func (c *Context) Abort() {
	c.index = abortIndex
}
*/
