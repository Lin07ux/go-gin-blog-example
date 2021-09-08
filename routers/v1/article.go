package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/app"
	"github.com/lin07ux/go-gin-example/pkg/e"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"github.com/lin07ux/go-gin-example/pkg/util"
	"github.com/lin07ux/go-gin-example/services"
	"github.com/unknwon/com"
	"net/http"
)

// 获取文章列表
func GetArticles(c *gin.Context) {
	maps := make(map[string]interface{})
	valid := validation.Validation{}
	response := app.Response{C:c}

	if arg := c.Query("state"); arg != "" {
		state := com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只能为 0、1")
	}

	if arg := c.Query("tag_id"); arg != "" {
		tagId := com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 1, "tag_id").Message("标签 ID 必须大于 0")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, valid.Errors[0].Message, nil)
		return
	}

	response.Send(e.Success, "", map[string]interface{}{
		"lists": models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps),
		"total": models.GetArticleTotal(maps),
	})
}

// 获取单个文章
func GetArticle(c *gin.Context) {
	response := app.Response{C: c}
	articleService := services.Article{ID: com.StrTo(c.Param("id")).MustInt()}

	valid := validation.Validation{}
	valid.Min(articleService.ID, 1, "id").Message("文章 ID 不存在")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, valid.Errors[0].Message, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		response.SetStatus(http.StatusInternalServerError).Send(e.ErrorGetArticleFail, "", nil)
	} else if article == nil {
		response.SetStatus(http.StatusNotFound).Send(e.ErrorNotExistArticle, "", nil)
	} else {
		response.Send(e.Success, "", article)
	}
}

// @Summary 添加文章
// @Produce json
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param state body int false "State"
// @Param created_by body string true "CreatedBy"
// @Success 200 {object} gin.H "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	response := app.Response{C: c}
	article := models.Article{
		TagID:      com.StrTo(c.PostForm("tag_id")).MustInt(),
		State:      com.StrTo(c.DefaultPostForm("state", "0")).MustInt(),
		Title:      c.PostForm("title"),
		Cover:      c.PostForm("cover"),
		Desc:       c.PostForm("desc"),
		Content:    c.PostForm("content"),
		CreatedBy:  c.PostForm("created_by"),
	}

	if message := validateCreateArticleData(&article); message != "" {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, message, nil)
		return
	}

	if ! models.ExistTagById(article.TagID) {
		response.SetStatus(http.StatusNotFound).Send(e.ErrorNotExistTag, "", nil)
		return
	}

	if id := models.AddArticle(&article); id > 0 {
		response.Send(e.Success, "", map[string]int{"id": id})
	} else {
		response.SetStatus(http.StatusExpectationFailed).Send(e.Error, "", nil)
	}
}

// @Summary 编辑文章
// @Produce json
// @Param id path int true "ID"
// @Param title body string false "Name"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedAt"
// @Success 200 {object} gin.H "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	response := app.Response{C: c}
	article := models.Article{
		TagID:      com.StrTo(c.PostForm("tag_id")).MustInt(),
		Title:      c.PostForm("title"),
		Cover:      c.PostForm("cover"),
		Desc:       c.PostForm("desc"),
		Content:    c.PostForm("content"),
		ModifiedBy: c.PostForm("modified_by"),
		State:      com.StrTo(c.DefaultPostForm("state", "-1")).MustInt(),
	}

	if message := validateUpdateArticleData(id, &article); message != "" {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, message, nil)
		return
	}

	if ! models.ExistArticleById(id) {
		response.SetStatus(http.StatusNotFound).Send(e.ErrorNotExistArticle, "", nil)
		return
	}

	if ! models.ExistTagById(article.TagID) {
		response.SetStatus(http.StatusNotFound).Send(e.ErrorNotExistTag, "", nil)
		return
	}

	if models.EditArticle(id, &article) {
		response.Send(e.Success, "", nil)
	} else {
		response.SetStatus(http.StatusExpectationFailed).Send(e.Error, "", nil)
	}
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	code := e.InvalidParams
	msg := ""

	if id > 0 {
		if models.ExistArticleById(id) {
			models.DeleteArticle(id)
			code = e.Success
		} else {
			code = e.ErrorNotExistArticle
		}

		msg = e.GetMsg(code)
	} else {
		msg = "文章 ID 必须大于 0"
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg" : msg,
		"data": make(map[string]string),
	})
}

// 创建文章数据校验
func validateCreateArticleData(article *models.Article) string {
	valid := validation.Validation{}
	valid.Min(article.TagID, 1, "tag_id").Message("标签 ID 必须大于 0")
	valid.Required(article.Title, "title").Message("文章标题不能为空")
	valid.MaxSize(article.Title, 100, "title").Message("文章标题不能超过 100 个字符")
	valid.Required(article.Cover, "cover").Message("文章封面图片地址不能为空")
	valid.MaxSize(article.Cover, 255, "cover").Message("文章封面图片地址不能超过 255 个字符")
	valid.Required(article.Desc, "desc").Message("文章简述不能为空")
	valid.MaxSize(article.Desc, 255, "desc").Message("文章简述不能超过 255 个字符")
	valid.Required(article.Content, "content").Message("文章内容不能为空")
	valid.MaxSize(article.Content, 65535, "content").Message("内容最长为 65535 个字符")
	valid.Required(article.CreatedBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(article.CreatedBy, 100, "created_by").Message("创建人最长为 100 个字符")
	valid.Range(article.State, 0, 1, "state").Message("文章状态只能为 0、1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		return valid.Errors[0].Message
	}

	return ""
}

// 更新文章数据校验
func validateUpdateArticleData(id int, article *models.Article) string {
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章 ID 必须大于 0")
	valid.MaxSize(article.Title, 100, "title").Message("标题最长为 100 个字符")
	valid.MaxSize(article.Cover, 255, "cover").Message("封面图片地址最长为 255 个字符")
	valid.MaxSize(article.Desc, 255, "desc").Message("简述最长为 255 个字符")
	valid.MaxSize(article.Content, 65535, "content").Message("内容最长为 65535 个字符")
	valid.Required(article.ModifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(article.ModifiedBy, 100, "modified_by").Message("修改人最长为 100 个字符")

	if article.State >= 0 {
		valid.Range(article.State, 0, 1, "state").Message("文章状态只能为 0、1")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		return valid.Errors[0].Message
	}

	return ""
}
