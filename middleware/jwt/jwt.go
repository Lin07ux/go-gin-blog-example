package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/pkg/e"
	"github.com/lin07ux/go-gin-example/pkg/util"
	"net/http"
	"time"
)

// jwt 中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.Success
		token := getAuthToken(c)

		if token == "" {
			code = e.ErrorAuthToken
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthCheckTokenFail
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ErrorAuthCheckTokenTimeout
			}
		}

		if code != e.Success {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg" : e.GetMsg(code),
				"data": nil,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

// 获取用户 token
func getAuthToken(c *gin.Context) (token string) {
	token = c.GetHeader("token")
	if token != "" {
		return
	}

	return c.Query("token")
}
