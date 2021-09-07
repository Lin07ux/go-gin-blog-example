package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/app"
	"github.com/lin07ux/go-gin-example/pkg/e"
	"github.com/lin07ux/go-gin-example/pkg/util"
	"net/http"
)

func GetAuth(c *gin.Context) {
	response := app.Response{C:c}
	username := c.PostForm("username")
	password := c.PostForm("password")

	if message := validateCredentials(username, password); message != "" {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, message, nil)
		return
	}

	if ! models.CheckAuth(username, password) {
		response.SetStatus(http.StatusForbidden).Send(e.ErrorAuth, "", nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		response.SetStatus(http.StatusInternalServerError).Send(e.ErrorAuthTokenGenerate, "", nil)
	} else {
		response.Send(e.Success, "", map[string]string{"token": token})
	}
}

// 登录凭证基本校验
func validateCredentials(username, password string) string {
	valid := validation.Validation{}
	valid.Required(username, "username").Message("账户名称不能为空")
	valid.MaxSize(username, 50, "username").Message("账户名称最长为 50 个字符")
	valid.Required(password, "password").Message("登录密码不能为空")
	valid.MaxSize(password, 50, "password").Message("登录密码最长为 50 个字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		return valid.Errors[0].Message
	}

	return ""
}
