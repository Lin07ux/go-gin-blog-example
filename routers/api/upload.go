package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/pkg/e"
	"github.com/lin07ux/go-gin-example/pkg/logging"
	"github.com/lin07ux/go-gin-example/pkg/upload"
	"net/http"
)

func UploadImage(c *gin.Context) {
	code := e.Success
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		code = e.Error
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg" : e.GetMsg(code),
			"data": data,
		})
		return
	}

	if image != nil {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()
		src := fullPath + imageName

		if ! upload.CheckImageExt(imageName) || ! upload.CheckImageSize(file) {
			code = e.ErrorUploadCheckImageFormat
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(err)
				code = e.ErrorUploadCheckImageBasic
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				logging.Warn(err)
				code = e.ErrorUploadSaveImageFail
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	} else {
		code = e.InvalidParams
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg" : e.GetMsg(code),
		"data": data,
	})
}
