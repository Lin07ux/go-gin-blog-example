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

// 获取文章列表
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

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

	code := e.InvalidParams
	msg := ""

	if ! valid.HasErrors() {
		code = e.Success
		msg = e.GetMsg(code)

		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
	} else {
		msg = valid.Errors[0].Message
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg" : msg,
		"data": data,
	})
}

// 获取单个文章
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章 ID 不存在")

	var data interface{}
	code := e.InvalidParams
	msg := ""

	if ! valid.HasErrors() {
		if models.ExistArticleById(id) {
			data = models.GetArticle(id)
			code = e.Success
		} else {
			code = e.ErrorNotExistArticle
		}
		msg = e.GetMsg(code)
	} else {
		msg = valid.Errors[0].Message
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg" : msg,
		"data": data,
	})
}

// 添加文章
func AddArticle(c *gin.Context) {
	article := models.Article{
		TagID:      com.StrTo(c.PostForm("tag_id")).MustInt(),
		State:      com.StrTo(c.DefaultPostForm("state", "0")).MustInt(),
		Title:      c.PostForm("title"),
		Desc:       c.PostForm("desc"),
		Content:    c.PostForm("content"),
		CreatedBy:  c.PostForm("created_by"),
	}

	code := e.InvalidParams
	result, msg := validateCreateArticleData(&article)

	if result {
		if models.ExistTagById(article.TagID) {
			models.AddArticle(&article)
			code = e.Success
		} else {
			code = e.ErrorNotExistTag
		}

		msg = e.GetMsg(code)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg" : msg,
		"data": make(map[string]interface{}),
	})
}

// 编辑文章
func EditArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	article := models.Article{
		TagID:      com.StrTo(c.PostForm("tag_id")).MustInt(),
		Title:      c.PostForm("title"),
		Desc:       c.PostForm("desc"),
		Content:    c.PostForm("content"),
		ModifiedBy: c.PostForm("modified_by"),
		State:      com.StrTo(c.DefaultPostForm("state", "-1")).MustInt(),
	}

	code := e.InvalidParams
	result, msg := validateUpdateArticleData(id, &article)

	if result {
		if models.ExistArticleById(id) {
			models.EditArticle(id, &article)
			code = e.Success
		} else {
			code = e.ErrorNotExistArticle
		}

		msg = e.GetMsg(code)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg" : msg,
		"data": make(map[string]string),
	})
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
func validateCreateArticleData(article *models.Article) (bool, string) {
	valid := validation.Validation{}
	valid.Min(article.TagID, 1, "tag_id").Message("标签 ID 必须大于 0")
	valid.Required(article.Title, "title").Message("文章标题不能为空")
	valid.MaxSize(article.Title, 100, "title").Message("文章标题不能超过 100 个字符")
	valid.Required(article.Desc, "desc").Message("文章简述不能为空")
	valid.MaxSize(article.Desc, 255, "desc").Message("文章简述不能超过 255 个字符")
	valid.Required(article.Content, "content").Message("文章内容不能为空")
	valid.MaxSize(article.Content, 65535, "content").Message("内容最长为 65535 个字符")
	valid.Required(article.CreatedBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(article.CreatedBy, 100, "created_by").Message("创建人最长为 100 个字符")
	valid.Range(article.State, 0, 1, "state").Message("文章状态只能为 0、1")

	if valid.HasErrors() {
		return false, valid.Errors[0].Message
	}

	return true, ""
}

// 更新文章数据校验
func validateUpdateArticleData(id int, article *models.Article) (bool, string) {
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章 ID 必须大于 0")
	valid.MaxSize(article.Title, 100, "title").Message("标题最长为 100 个字符")
	valid.MaxSize(article.Desc, 255, "desc").Message("简述最长为 255 个字符")
	valid.MaxSize(article.Content, 65535, "content").Message("内容最长为 65535 个字符")
	valid.Required(article.ModifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(article.ModifiedBy, 100, "modified_by").Message("修改人最长为 100 个字符")

	if article.State >= 0 {
		valid.Range(article.State, 0, 1, "state").Message("文章状态只能为 0、1")
	}

	if valid.HasErrors() {
		return false, valid.Errors[0].Message
	}

	return true, ""
}
