package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/pkg/app"
	"github.com/lin07ux/go-gin-example/pkg/e"
	"github.com/lin07ux/go-gin-example/pkg/logging"
	"github.com/lin07ux/go-gin-example/pkg/upload"
	"net/http"
)

func UploadImage(c *gin.Context) {
	response := app.Response{C: c}

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		response.SetStatus(http.StatusInternalServerError).Send(e.Error, "", nil)
		return
	}

	if image == nil {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, "", nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + imageName

	if ! upload.CheckImageExt(imageName) || ! upload.CheckImageSize(file) {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.ErrorUploadCheckImageFormat, "", nil)
		return
	}

	if err := upload.CheckImage(fullPath); err != nil {
		logging.Warn(err)
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.ErrorUploadCheckImageBasic, "", nil)
		return
	}

	if err := c.SaveUploadedFile(image, src); err != nil {
		logging.Warn(err)
		response.SetStatus(http.StatusInternalServerError).Send(e.ErrorUploadSaveImageFail, "", nil)
	} else {
		response.Send(e.Success, "", map[string]string{
			"image_url":      upload.GetImageFullUrl(imageName),
			"image_save_url": savePath + imageName,
		})
	}
}
