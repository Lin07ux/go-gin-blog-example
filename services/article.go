package services

import (
	"github.com/lin07ux/go-gin-example/cache"
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/logging"
)

type Article struct {
	ID    int
	TagID int
	State int

	Title      string
	Cover      string
	Desc       string
	Content    string
	CreatedBy  string
	ModifiedBy string

	PageNum int
	PageSize int
}

// 文章详情
func (a *Article) Detail() (*models.Article, error) {
	articleCache := cache.Article{ID: a.ID}

	if article := articleCache.GetDetail(); article != nil {
		return article, nil
	}

	article, err := models.GetArticle(a.ID)

	if err != nil {
		logging.Warn(err)
	} else if article != nil {
		articleCache.SetDetail(article)
	}

	return article, err
}

// 文章列表
func (a *Article) List() ([]*models.Article, error) {
	articleCache := cache.Article{
		TagID:    a.TagID,
		State:    a.State,
		PageNum:  a.PageSize,
		PageSize: a.PageSize,
	}

	if articles := articleCache.GetList(); articles != nil {
		return articles, nil
	}

	articles, err := models.GetArticles(a.getListOffset(), a.PageSize, a.getQueryMaps())
	if err != nil {
		logging.Warn(err)
	} else if articles != nil {
		articleCache.SetList(articles)
	}

	return articles, err
}

// 文章数量
func (a *Article) Count() (total int64, err error) {
	total, err = models.GetArticleTotal(a.getQueryMaps())
	if err != nil {
		logging.Warn(err)
	}

	return
}

// 增加文章
func (a *Article) Create() int {
	var err error
	a.ID, err = models.AddArticle(&models.Article{
		TagID:      a.TagID,
		Title:      a.Title,
		Cover:      a.Cover,
		Desc:       a.Desc,
		Content:    a.Content,
		CreatedBy:  a.CreatedBy,
		State:      a.State,
	})

	if err != nil {
		logging.Warn(err)
		return 0
	}

	return a.ID
}

// 更新文章
func (a *Article) Update() bool {
	data := make(map[string]interface{})
	data["modified_by"] = a.ModifiedBy
	if a.TagID > 0 {
		data["tag_id"] = a.TagID
	}
	if a.Title != "" {
		data["title"] = a.Title
	}
	if a.Cover != "" {
		data["cover"] = a.Cover
	}
	if a.Desc != "" {
		data["desc"] = a.Desc
	}
	if a.Content != "" {
		data["content"] = a.Content
	}
	if a.State >= 0 {
		data["state"] = a.State
	}

	if err := models.EditArticle(a.ID, data); err != nil {
		logging.Warn(err)
		return false
	}

	// 清除缓存
	articleArticle := cache.Article{ID: a.ID}
	articleArticle.CleanForArticle()

	return true
}

// 删除文章
func (a *Article) Delete() bool {
	if err := models.DeleteArticle(a.ID); err != nil {
		logging.Warn(err)
		return false
	}

	// 清除缓存
	articleArticle := cache.Article{ID: a.ID}
	articleArticle.CleanForArticle()

	return true
}

// 判断文章是否存在
func (a *Article) ExistsById() bool {
	exists, err := models.ExistArticleById(a.ID)
	if err != nil {
		logging.Warn(err)
		return false
	}

	return exists
}

// 获取文章列表的查询参数
func (a *Article) getQueryMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if a.State >= 0 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}

	return maps
}

// 获取文章列表的分页起始位置
func (a *Article) getListOffset() int {
	return (a.PageNum - 1) * a.PageSize
}
