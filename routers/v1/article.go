package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/app"
	"github.com/lin07ux/go-gin-example/pkg/e"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"github.com/lin07ux/go-gin-example/services"
	"github.com/unknwon/com"
	"net/http"
)

// @Summary 获取文章列表
// @Produce json
// @param state query int false "文章状态" Enums(0, 1)
// @Param tag_id query int false "所属标签 ID" minimum(1)
// @Param page query int false "分页页码" minimum(1) default(1)
// @Success 200 {object} app.ResponseBody{data={list=[]models.Article,total:int}} "ok"
// @Failure 422 {object} app.ResponseBody "请求参数错误"
// @Failure 500 {object} app.ResponseBody "获取文章数据失败"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	response := app.Response{C: c}
	articleService := services.Article{
		TagID:    com.StrTo(c.DefaultQuery("state", "-1")).MustInt(),
		State:    com.StrTo(c.DefaultQuery("tag_id", "-1")).MustInt(),
		PageNum:  com.StrTo(c.Query("page")).MustInt(),
		PageSize: setting.AppSetting.PageSize,
	}

	if message := validateArticlesQueries(&articleService); message != "" {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, message, nil)
		return
	}

	total, err := articleService.Count()
	if err != nil {
		response.SetStatus(http.StatusInternalServerError).Send(e.ErrorCountArticlesFail, "", nil)
		return
	}

	var articles = make([]*models.Article, 0, 1)
	if total > 0 {
		articles, err = articleService.List()
		if err != nil {
			response.SetStatus(http.StatusInternalServerError).Send(e.ErrorGetArticlesFail, "", nil)
			return
		}
	}

	response.Send(e.Success, "", map[string]interface{}{
		"lists": articles,
		"total": total,
	})
}

// @Summary 获取单个文章
// @Produce json
// @param id path int true "文章 ID" minimum(1)
// @Success 200 {object} app.ResponseBody{data=models.Article} "ok"
// @Failure 404 {object} app.ResponseBody{code=int(404)} "文章不存在"
// @Failure 422 {object} app.ResponseBody{code=int(422)} "请求参数错误"
// @Failure 500 {object} app.ResponseBody "获取文章失败"
// @Router /api/v1/article/{id} [get]
func GetArticle(c *gin.Context) {
	response := app.Response{C: c}
	articleService := services.Article{ID: com.StrTo(c.Param("id")).MustInt()}

	if message := validateArticleId(&articleService); message != "" {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, message, nil)
		return
	}

	if article, err := articleService.Detail(); err != nil {
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
// @Success 200 {object} app.ResponseBody "ok"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	response := app.Response{C: c}
	articleService := services.Article{
		TagID:      com.StrTo(c.PostForm("tag_id")).MustInt(),
		State:      com.StrTo(c.DefaultPostForm("state", "0")).MustInt(),
		Title:      c.PostForm("title"),
		Cover:      c.PostForm("cover"),
		Desc:       c.PostForm("desc"),
		Content:    c.PostForm("content"),
		CreatedBy:  c.PostForm("created_by"),
	}

	if message := validateCreateArticleData(&articleService); message != "" {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, message, nil)
		return
	}

	if ! models.ExistTagById(articleService.TagID) {
		response.SetStatus(http.StatusNotFound).Send(e.ErrorNotExistTag, "", nil)
		return
	}

	if articleService.Create() <= 0 {
		response.SetStatus(http.StatusExpectationFailed).Send(e.Error, "", nil)
	} else {
		response.Send(e.Success, "", map[string]int{"id": articleService.ID})
	}
}

// @Summary 编辑文章
// @Produce json
// @Param id path int true "ID"
// @Param title body string false "Name"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedAt"
// @Success 200 {object} app.ResponseBody "ok"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	response := app.Response{C: c}
	article := services.Article{
		ID:         com.StrTo(c.Param("id")).MustInt(),
		TagID:      com.StrTo(c.DefaultPostForm("tag_id", "0")).MustInt(),
		State:      com.StrTo(c.DefaultPostForm("state", "-1")).MustInt(),
		Title:      c.PostForm("title"),
		Cover:      c.PostForm("cover"),
		Desc:       c.PostForm("desc"),
		Content:    c.PostForm("content"),
		ModifiedBy: c.PostForm("modified_by"),
	}

	if message := validateUpdateArticleData(&article); message != "" {
		response.SetStatus(http.StatusUnprocessableEntity).Send(e.InvalidParams, message, nil)
		return
	}

	if ! article.ExistsById() {
		response.SetStatus(http.StatusNotFound).Send(e.ErrorNotExistArticle, "", nil)
		return
	}

	if ! models.ExistTagById(article.TagID) {
		response.SetStatus(http.StatusNotFound).Send(e.ErrorNotExistTag, "", nil)
		return
	}

	if ! article.Update() {
		response.SetStatus(http.StatusExpectationFailed).Send(e.Error, "", nil)
	} else {
		response.Send(e.Success, "", nil)
	}
}

// @Summary 删除文章
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} app.ResponseBody "ok"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	response := app.Response{C:c}
	articleService := services.Article{ID: com.StrTo(c.Param("id")).MustInt()}

	if message := validateArticleId(&articleService); message != "" {
		response.SetStatus(http.StatusNotFound).Send(e.InvalidParams, message, nil)
		return
	}

	if ! articleService.ExistsById() {
		response.SetStatus(http.StatusNotFound).Send(e.ErrorNotExistArticle, "", nil)
		return
	}

	if ! articleService.Delete() {
		response.SetStatus(http.StatusInternalServerError).Send(e.ErrorDeleteArticleFail, "", nil)
	} else {
		response.Send(e.Success, "", nil)
	}
}

// 验证文章 ID
func validateArticleId(article *services.Article) string {
	valid := validation.Validation{}
	valid.Min(article.ID, 1, "id").Message("文章 ID 不存在")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		return valid.Errors[0].Message
	}

	return ""
}

// 验证文章列表查询参数
func validateArticlesQueries(article *services.Article) string {
	valid := validation.Validation{}

	if article.State >= 0 {
		valid.Range(article.State, 0, 1, "state").Message("状态只能为 0、1")
	}
	if article.TagID > 0 {
		valid.Min(article.TagID, 1, "tag_id").Message("标签 ID 必须大于 0")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		return valid.Errors[0].Message
	}

	return ""
}

// 创建文章数据校验
func validateCreateArticleData(article *services.Article) string {
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
func validateUpdateArticleData(article *services.Article) string {
	valid := validation.Validation{}
	valid.Min(article.ID, 1, "id").Message("文章 ID 必须大于 0")
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
