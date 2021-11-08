package models

import (
	"github.com/dgrijalva/jwt-go"
)

// ID\Nickname\AuthorityId是payload中的主要内容
// 在加密的时候要用到CustomClaims，解密的时候也需要用到CustomClaims
type CustomClaims struct {
	ID                 int
	Nickname           string
	AuthorityId        uint //是role，有些接口只能给后端系统的管理系统才能使用
	jwt.StandardClaims      // 是一些标准的内容
}

/**
type StandardClaims struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"` //过期时间
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}
*/
