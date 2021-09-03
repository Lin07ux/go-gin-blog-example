package e

var MsgFlags = map[int]string {
	Success:                    "ok",
	Error:                      "fail",
	InvalidParams:              "请求参数错误",
	ErrorExistTag:              "已存在该标签名称",
	ErrorNotExistTag:           "该标签不存在",
	ErrorNotExistArticle:       "该文章不存在",
	ErrorAuthCheckTokenFail:    "Token 鉴权失败",
	ErrorAuthCheckTokenTimeout: "Token 已超时",
	ErrorAuthTokenGenerate:     "Token 生成失败",
	ErrorAuthToken:             "Token 鉴权失败",
	ErrorAuth:                  "登录信息不匹配",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[Error]
}