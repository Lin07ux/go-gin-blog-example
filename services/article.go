package services

import (
	"github.com/lin07ux/go-gin-example/cache"
	"github.com/lin07ux/go-gin-example/models"
	"log"
)

type Article struct {
	ID int
}

func (a *Article) Get() (*models.Article, error) {
	cacheArticle := cache.Article{ID: a.ID}

	if article := cacheArticle.GetDetail(); article != nil {
		log.Println("get article detail fro cache")
		return article, nil
	}

	article, err := models.GetArticle(a.ID)
	if article != nil && err == nil {
		cacheArticle.SetDetail(article)
	}

	return article, err
}

func (a *Article) ExistsById() bool {
	return models.ExistArticleById(a.ID)
}
