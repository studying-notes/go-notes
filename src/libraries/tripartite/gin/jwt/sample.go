package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// token 存储信息结构体
type claims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

// token 过期时间
const tokenExpireDuration = time.Minute * 3

// 密钥
var secret = []byte("secret")

func GenToken(user string) (string, error) {
	c := claims{"user", jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
		Issuer:    "rustle",
	}}

	// 指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 指定 secret 签名并获得完整的编码后的字符串
	return token.SignedString(secret)
}

func ParseToke(signedString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(signedString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func main() {
	user := "example"
	token, _ := GenToken(user)
	fmt.Println(token)

	parsed, _ := ParseToke(token)
	fmt.Printf("%+v\n", *parsed) // 结构体的成员名称和值
}
