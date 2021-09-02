package v1

import (
	"github.com/astaxie/beego/validation"
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
	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	createdBy := c.PostForm("created_by")

	code := e.InvalidParams
	msg := ""

	result, message := validateTagData(name, state, createdBy)
	if result {
		if ! models.ExistTagByName(name) {
			code = e.Success
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ErrorExistTag
		}
		msg = e.GetMsg(code)
	} else {
		msg = message
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": msg,
		"data": make(map[string]string),
	})
}

// 修改文章标签
func EditTag(c *gin.Context) {
	//
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	//
}

func validateTagData(name string, state int, createdBy string) (result bool, message string) {
	valid := validation.Validation{}
	valid.Required(name, "name").Message("标签名称不能为空")
	valid.MaxSize(name, 100, "name").Message("标签名称最长为 100 个字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为 100 个字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许为 0、1")

	if valid.HasErrors() {
		return false, valid.Errors[0].Message
	}

	return true, ""
}
