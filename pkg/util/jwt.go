package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"time"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 生成 token
func GenerateToken(username, password string) (string, error) {
	claims := Claims{
		Username:       username,
		Password:       password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt:  time.Now().Add(3 * time.Hour).Unix(),
			Issuer:     "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokenClaims.SignedString(jwtSecret)
}

// 解析 token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}