package main

import (
	"context"
	"golang_learn/golang_learn/models/proto"
	validatorTrans "golang_learn/golang_learn/models/validator"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

/**
  传统的非前后端分离的系统中，使用的cookie和session的方式！如果在微服务中使用这种方式:
微服务是互相隔离的，数据库也是独立的。可以将session保存在redis集群中。

为什么json web token机制适合在微服务中应用呢？他是通过加密的技术完成的，为什么加密可以完成这个效果呢？
登录的原理：给浏览器一个加密的字符串，这个字符串是不能够随意伪造的，只要达到这一点就能够达到 认证 的作用。只有微服务可以解密，浏览器是不能够解密的。

jwt是一个开放标准，用于作为json对象在各方之间安全地传输信息，可以被验证和信任的，因为他是数字签名的。

什么时候使用jwt呢？
1、授权
这是使用jwt的最常见的场景。一旦用户登录，后续的每一个请求都包含jwt，允许用户访问该令牌允许的路由、服务和资源。单点登录是现在最广泛使用jwt的一个特性，因为他的开销很小，并且可以轻松的跨域使用
2、信息交换
对于安全的在各方传输信息而言，jwt无疑是一种很好的方式。因为jwt可以被签名。例如，用公钥/私钥，可以确定发送者就是他们所说的哪个人。另外，由于签名是使用头和有效负载计算的，还可以验证内容有没有被篡改

jwt的结构：它有header、payload、signature构成
header.payload.signature
aaaaaa.bbbbbbb.cccccccc

header：典型的由两部分组成：token的类型，即jwt。和算法名称组成，如rsa、hmac、sha256
payload：它包含声明。声明式关于实体（通常是用户）和其他数据的声明。 声明有三种类型：Registered、Public、Private

	Registered：这里有一组预定义的声明，他们不是强制的，但是推荐。比如iss（issuer），exp（expiration time），sub（subject），aud（audience）等
	Public：可以随意定义
	Private：用于在同意使用他们的各方之间共享信息，并且不是注册的或公开的声明

下面是一个例子：
{
	"sub":"12345677",
	"name":"json",
	"admin":true
}
注意：对payload进行base64编码就得到jwt的第二部分。不要在jwt的payload或header中防止敏感信息，除非他们是加密的。

signature：为了得到这部分，必须有编码过的header、编码过的payload、一个秘钥，签名算法是header中指定的那个，然后对他们签名即可
例如：HMACSHA256(base64UrlEncode(header)+"."+base64UrlEncode(payload),secret)
签名是用于验证，消息在传递的过程中有没有被变更，并且对于使用私钥签名的token，还可以验证jwt的发送方式否为他所有成的发送方。

通过加密和解密验证真实性，所以就不需要在服务端存储了！
jwt.io
*/

/**
jwt是如何工作的？
当用户登录以后，会返回一个jwt。此后，jwt就是用户凭证了，你必须非常小心防止出现安全问题。一般而言，你保存令牌的时候不应该超过你锁需要他的时间。
无论何时用户想要访问受保护的路由或者资源的时候，用户代理（通常是浏览器）都应该带上jwt，典型的，通常放在Authorization header中，用Bearer schema。
header应该看起来是这样的：
Authorization：Bearer

	服务端上的受保护的路由将会检查Authorization header中的jwt是否有效，如果有效，则用户可以访问受保护的资源，如果jwt包含足够多的必须的数据，
那么就可以减少对某些操作的数据库查询的需求，尽管可能并不总是如此。
	如果token是在Authorization header中发送的，那么跨源资源共享（CORS）将不会成为问题，因为它不使用cookie

操作的流程：
1、应用或者客户端向授权服务器，请求授权。例如，如果用授权码流程的话，就是/oauth/authorize
2、当授权被许可以后，授权服务器会返回一个access token给应用
3、应用使用access token访问受保护的资源，如api

一般token放在header中
*/

/**
基于token的身份认证 与 基于服务器的身份认证
1、基于服务器的身份认证
在讨论基于token的身份认证是如何工作的以及他的好处之前，我们先来看一下以前我们是怎么做的：
HTTP协议是无状态的，也就是说没如果饿哦们已经认证了一个用户，那么下一次请求的时候，服务器不知道这个用户，必须再次认证
传统的做法是将已经认证过的用户信息存储在服务器上，比如session。用户瞎子啊请求的时候带着session id，然后服务器以此检查用户是否验证过。
这种基于服务器的身份认证方式存在一些问题：
a、sessions：每次认证通过以后，服务器需要创建一条记录保存用户信息，通常是在内存中，随着认证通过的用户越来越多，服务器在保存session的开销会越来越大
b、scalability：由于session是在内存中，这就带来一些扩展性的问题
c、CORS：当我们想要扩展应用的时候，让我们的数据被多个移动设备使用时，我们必须考虑跨资源共享问题。当使用ajax调用从另一个域名下，获取资源时，我们可能会遇到禁止请求的问题。
d、CSRF：用户很容易收到CSRF攻击

jwt与session的差异：
相同点：他们都是存储用户信息，session存储在服务端，jwt存储在客户端
session方式存储用户信息的最大问题在于要占用大量服务器内存，增加了服务器的开销；而jwt将用户状态分散到客户端中，可以明显减轻服务端的内存的压力。

2、基于token的身份认证是如何工作的
基于token的身份认证是无状态的，服务器或session中不会存储任何用户信息。
没有会话信息意味着应用程序可以根据需要扩展和添加更多的机器，而不必担心用户登录的位置。

主流程如下：
a、用户携带用户名和密码请求访问
b、服务器校验用户凭证
c、应用提供一个token给客户端
d、客户端存储token，并且在虽有的每一个请求中都携带这个它
e、服务器校验token并返回数据

注意：
a、每一次请求都需要token
b、token应在放在请求header中
c、我们还需要将服务器设置为接受来自手游域的请求，用access-control-allow-origin：*

用token的好处：
1、无状态和扩展性
tokens存储在客户端。完全没有状态，可扩展。我们的负载均衡器可以将用户传递到任意服务器，因为在任何地方都没有状态或会话信息。

2、安全
token不是cookie。每次请求的时候token都会被发送。而且，由于没有cookie被发送，还有助于防止CSRF攻击。即使在你的实现中token存储在客户端的cookie中，这个cookie也只是一个存储机制，
而非身份认证机制。没有基于会话的信息可以操作，因为没有会话。
token在一段时间以后会过期，这个时候用户需要重新登录。这有助于我们保持安全。还有一个概念叫token撤销，它允许我们根据相同的授权许可使特定的token，甚至一组token无效。

jwt与oauth的区别：
a、oauth2是一个授权框架，jwt是一个认证协议
b、无论使用哪种方式切记用https来保证数据的安全性
c、oauth2用子啊使用第三方账号登录的情况（比如使用微博、qq、github登录某个app），而jwt是用在前后端分离，需要简单的对后台api进行保护时使用！！！
*/
// 这里需要设置

func main() {
	// 我们需要的就是验证用户
	//生成token

}

type PassWordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻， 自定义validator
	PassWord  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	//TODO：err断言为validator.ValidationErrors失败
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(validatorTrans.Trans)), //Translate 翻译
	})
}

// user
func PasswordLogin(c *gin.Context) {
	passwordLoginLoginForm := PassWordLoginForm{}
	//表单验证:对passwordLoginForm执行验证
	if err := c.ShouldBind(&passwordLoginLoginForm); err != nil {
		//表单验证不通过
		HandleValidatorError(c, err)
	}
	//验证码
	if store.Verify(passwordLoginLoginForm.CaptchaId, passwordLoginLoginForm.Captcha, false) {
		c.JSON(http.StatusBadRequest, gin.H{
			"capture": "验证码错误",
		})
	}
	//登录逻辑
	// 调用rpc 的user service 端接口，根据手机号获取用户信息
	if _, err := userServiceClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginLoginForm.Mobile,
	}); err != nil {
		//判断返回的状态码
		if s, ok := status.FromError(err); ok {
			/**
			const (
				// OK is returned on success.
				OK Code = 0
				Canceled Code = 1
				Unknown Code = 2
				InvalidArgument Code = 3
				DeadlineExceeded Code = 4
				NotFound Code = 5
				AlreadyExists Code = 6
				PermissionDenied Code = 7
				ResourceExhausted Code = 8
				FailedPrecondition Code = 9
				Aborted Code = 10
				OutOfRange Code = 11
				Unimplemented Code = 12
				Internal Code = 13
				Unavailable Code = 14
				DataLoss Code = 15
				Unauthenticated Code = 16
				_maxCode = 17
			)
			*/
			switch s.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}

		}
	} else {
		// 校验获得到的用户密码

	}

}

var userServiceClient proto.UserClient

// 这里是
func GetUserByMobile(ctx context.Context, mobileRequest *proto.MobileRequest) {
	//根据
}
