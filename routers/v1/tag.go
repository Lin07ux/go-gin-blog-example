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

	result, message := validateCreateTagData(name, state, createdBy)
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
	id := com.StrTo(c.Param("id")).MustInt()
	state := com.StrTo(c.DefaultPostForm("state", "-1")).MustInt()
	name := c.PostForm("name")
	modifiedBy := c.PostForm("modified_by")

	code := e.InvalidParams
	msg := ""

	result, message := validateUpdateTagData(id, state, name, modifiedBy)
	if result {
		if models.ExistTagById(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state >= 0 {
				data["state"] = state
			}
			models.EditTag(id, data)
			code = e.Success
		} else {
			code = e.ErrorNotExistTag
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

// 删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("标签 ID 不存在")

	code := e.Success
	msg := ""

	if ! valid.HasErrors() {
		if models.ExistTagById(id) {
			models.DeleteTag(id)
		} else {
			code = e.ErrorNotExistTag
		}
		msg = e.GetMsg(code)
	} else {
		code = e.InvalidParams
		msg = valid.Errors[0].Message
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": msg,
		"data": make(map[string]string),
	})
}

// 校验创建文章标签的数据
func validateCreateTagData(name string, state int, createdBy string) (bool, string) {
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

// 校验更新文章标签的数据
func validateUpdateTagData(id int, state int, name string, modifiedBy string) (bool, string) {
	valid := validation.Validation{}
	valid.Required(id, "id").Message("标签 ID 不能为空")
	valid.Min(id, 1, "id").Message("标签 ID 不存在")
	valid.MaxSize(name, 100, "name").Message("标签名称最长为 100 个字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为 100 个字符")

	if state >= 0 {
		valid.Range(state, 0, 1, "state").Message("标签状态只能为 0、1")
	}

	if valid.HasErrors() {
		return false, valid.Errors[0].Message
	}

	return true, ""
}
