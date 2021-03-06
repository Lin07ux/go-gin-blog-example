package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/pkg/e"
	"net/http"
)

type Response struct {
	C      *gin.Context
	status int
}

type ResponseBody struct {
	Code int `json:"code" example:"0"`
	Message string `json:"message" example:"ok"`
	Data interface{} `json:"data"`
}

func (g *Response) SetStatus(code int) *Response {
	g.status = code

	return g
}

func (g *Response) Send(code int, message string, data interface{}) {
	httpCode := g.status
	if httpCode <= 0 {
		httpCode = http.StatusOK
	}

	if message == "" {
		message = e.GetMsg(code)
	}

	if data == nil {
		data = make(map[string]string)
	}

	g.C.JSON(httpCode, ResponseBody{
		Code:    code,
		Message: message,
		Data:    data,
	})
}