package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/e"
	"github.com/lin07ux/go-gin-example/pkg/util"
	"net/http"
)

func GetAuth(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	code := e.InvalidParams
	data := make(map[string]interface{})
	result, msg := validateCredentials(username, password)

	if result {
		if models.CheckAuth(username, password) {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ErrorAuthTokenGenerate
			} else {
				code = e.Success
				data["token"] = token
			}
		} else {
			code = e.ErrorAuth
		}
		msg = e.GetMsg(code)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg" : msg,
		"data": data,
	})
}

// 登录凭证基本校验
func validateCredentials(username, password string) (bool, string) {
	valid := validation.Validation{}
	valid.Required(username, "username").Message("账户名称不能为空")
	valid.MaxSize(username, 50, "username").Message("账户名称最长为 50 个字符")
	valid.Required(password, "password").Message("登录密码不能为空")
	valid.MaxSize(password, 50, "password").Message("登录密码最长为 50 个字符")

	if valid.HasErrors() {
		return false, valid.Errors[0].Message
	}

	return true, ""
}
