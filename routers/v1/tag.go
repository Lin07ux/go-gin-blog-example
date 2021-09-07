package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/app"
	"github.com/lin07ux/go-gin-example/pkg/e"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"github.com/lin07ux/go-gin-example/pkg/util"
	"github.com/unknwon/com"
	"net/http"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	response := app.Response{C:c}
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name := c.Query("name"); name != "" {
		maps["name"] = name
	}

	if state := c.Query("state"); state != "" {
		maps["state"] = com.StrTo(state).MustInt()
	}

	data["lists"] = models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	response.Send(e.Success, "", data)
}

// 新增文章标签
func AddTag(c *gin.Context) {
	response := app.Response{C:c}
	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	createdBy := c.PostForm("created_by")

	if message := validateCreateTagData(name, state, createdBy); message != "" {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, message, nil)
		return
	}

	if models.ExistTagByName(name) {
		response.SetStatus(http.StatusPreconditionFailed).Send(e.ErrorExistTag, "", nil)
	} else {
		id := models.AddTag(name, state, createdBy)
		response.Send(e.Success, "", map[string]int{"id": id})
	}
}

// 修改文章标签
func EditTag(c *gin.Context) {
	response := app.Response{C:c}
	id := com.StrTo(c.Param("id")).MustInt()
	tag := &models.Tag{
		Name:       c.PostForm("name"),
		State:      com.StrTo(c.DefaultPostForm("state", "-1")).MustInt(),
		ModifiedBy: c.PostForm("modified_by"),
	}

	if message := validateUpdateTagData(id, tag); message != "" {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, message, nil)
		return
	}

	tagModel := models.GetTagById(id)
	if tagModel.ID <= 0 {
		response.SetStatus(http.StatusPreconditionFailed).Send(e.ErrorNotExistTag, "", nil)
		return
	}

	if tagModel.Name != tag.Name && models.ExistTagByName(tag.Name) {
		response.SetStatus(http.StatusPreconditionFailed).Send(e.ErrorExistTag, "", nil)
		return
	}

	models.EditTag(id, tag)
	response.Send(e.Success, "", nil)
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	response := app.Response{C:c}

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("标签 ID 不存在")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, valid.Errors[0].Message, nil)
		return
	}

	if ! models.ExistTagById(id) {
		response.SetStatus(http.StatusNotFound).Send(e.ErrorNotExistTag, "", nil)
	} else {
		models.DeleteTag(id)
		response.Send(e.Success, "", nil)
	}
}

// 校验创建文章标签的数据
func validateCreateTagData(name string, state int, createdBy string) string {
	valid := validation.Validation{}
	valid.Required(name, "name").Message("标签名称不能为空")
	valid.MaxSize(name, 100, "name").Message("标签名称最长为 100 个字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为 100 个字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许为 0、1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		return valid.Errors[0].Message
	}

	return ""
}

// 校验更新文章标签的数据
func validateUpdateTagData(id int, tag *models.Tag) string {
	valid := validation.Validation{}
	valid.Required(id, "id").Message("标签 ID 不能为空")
	valid.Min(id, 1, "id").Message("标签 ID 不存在")
	valid.MaxSize(tag.Name, 100, "name").Message("标签名称最长为 100 个字符")
	valid.Required(tag.ModifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(tag.ModifiedBy, 100, "modified_by").Message("修改人最长为 100 个字符")

	if tag.State >= 0 {
		valid.Range(tag.State, 0, 1, "state").Message("标签状态只能为 0、1")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		return valid.Errors[0].Message
	}

	return ""
}
