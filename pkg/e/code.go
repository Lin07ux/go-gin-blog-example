package e

const (
	Success       = 200
	InvalidParams = 422
	Error         = 500

	ErrorExistTag        = 10001
	ErrorNotExistTag     = 10002
	ErrorNotExistArticle = 10003

	ErrorAuthCheckTokenFail    = 20001
	ErrorAuthCheckTokenTimeout = 20002
	ErrorAuthTokenGenerate     = 2003
	ErrorAuthToken             = 2004
	ErrorAuth                  = 2005

	ErrorUploadCheckImageBasic  = 30001
	ErrorUploadCheckImageFormat = 30002
	ErrorUploadSaveImageFail    = 3003
)
