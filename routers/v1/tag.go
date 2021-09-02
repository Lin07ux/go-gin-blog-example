package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/e"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"github.com/lin07ux/go-gin-example/pkg/util"
	"github.com/unknwon/com"
	"net/http"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name := c.Query("name"); name != "" {
		maps["name"] = name
	}

	if state := c.Query("state"); state != "" {
		maps["state"] = com.StrTo(state).MustInt()
	}

	code := e.Success
	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

// 新增文章标签
func AddTag(c *gin.Context) {
	//
}

// 修改文章标签
func EditTag(c *gin.Context) {
	//
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	//
}
