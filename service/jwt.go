package service

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("https://gitee.com/phpmysic/smash-golden-eggs")
var issuer = "127.0.0.1"
var subject = "user mobile"

type Claims struct {
	UserId string
	jwt.StandardClaims
}

// Setting 颁发token
func Setting(mobile string) (string,error) {
	expireTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserId: mobile,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,  // 签名颁发者
			Subject:   subject, //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Getting 解析token

func Getting(tokenString string) (*jwt.Token, string, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, Claims.UserId, err
}