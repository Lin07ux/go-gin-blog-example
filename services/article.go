package services

import (
	"github.com/lin07ux/go-gin-example/cache"
	"github.com/lin07ux/go-gin-example/models"
	"log"
)

type Article struct {
	ID    int
	TagID int
	State int

	PageNum int
	PageSize int
}

// 文章详情
func (a *Article) Detail() (*models.Article, error) {
	articleCache := cache.Article{ID: a.ID}

	if article := articleCache.GetDetail(); article != nil {
		log.Println("get article detail fro articleCache")
		return article, nil
	}

	article, err := models.GetArticle(a.ID)
	if article != nil && err == nil {
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

	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getQueryMaps())
	if articles != nil && err == nil {
		articleCache.SetList(articles)
	}

	return articles, err
}

// 文章数量
func (a *Article) Count() (int64, error) {
	return models.GetArticleTotal(a.getQueryMaps())
}

func (a *Article) ExistsById() bool {
	return models.ExistArticleById(a.ID)
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
