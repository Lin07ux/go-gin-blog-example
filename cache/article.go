package cache

import (
	"fmt"
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/gredis"
	"github.com/lin07ux/go-gin-example/pkg/logging"
	"log"
	"strconv"
	"strings"
	"time"
)

type Article struct {
	ID    int
	TagID int
	State int

	PageNum  int
	PageSize int
}

// GetDetail get article detail form cache
func (a *Article) GetDetail() *models.Article {
	key := a.getArticleDetailKey()

	exist, err := gredis.Exists(key)
	if ! exist || err != nil {
		if err != nil {
			logging.Warn(err)
		}
		return nil
	}

	detail := &models.Article{}
	if err := gredis.Get(key, detail); err != nil {
		logging.Warn(err)
		return nil
	}

	return detail
}

// SetDetail set article detail to cache
func (a *Article) SetDetail(article *models.Article) bool {
	err := gredis.Set(a.getArticleDetailKey(), article, 3600 * time.Second)
	if err != nil {
		logging.Warn(err)
		return false
	}

	return true
}

// DelDetail will remove article detail cache
func (a *Article) DelDetail() bool {
	success, err := gredis.Delete(a.getArticleDetailKey())
	if err != nil {
		logging.Warn(err)
		return false
	}

	return success
}

// GetList get articles list from cache
func (a *Article) GetList() []*models.Article {
	key := a.getArticleListKey()

	exist, err := gredis.Exists(key)
	if ! exist || err != nil {
		if err != nil {
			logging.Warn(err)
		}
		return nil
	}

	var articles []*models.Article
	if err := gredis.Get(a.getArticleListKey(), articles); err != nil {
		logging.Warn(err)
		return nil
	}

	return articles
}

// SetList set articles list to cache
func (a *Article) SetList(articles []*models.Article) bool {
	key := a.getArticleListKey()

	if err := gredis.Set(key, articles, 600 * time.Second); err != nil {
		logging.Warn(err)
		return false
	}

	return true
}

// CleanForArticle clean the article detail cache and all article list cache
func (a *Article) CleanForArticle() bool {
	success := a.DelDetail()

	state := a.State
	a.State = -1

	if err := gredis.LikeDeletes(a.getArticleListKey()); err != nil {
		log.Printf("clean list error: %v", err)
		logging.Warn(err)
		success = false
	}

	a.State = state

	return success
}

// getArticleDetailKey build the key for article detail cache
func (a *Article) getArticleDetailKey() string {
	return fmt.Sprintf("article:%d", a.ID)
}

// getArticleListKey build the key for articles list cache
func (a *Article) getArticleListKey() string {
	var keys []string

	if a.TagID > 0 {
		keys = append(keys, strconv.Itoa(a.TagID))
	}
	if a.State >= 0 {
		keys = append(keys, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		keys = append(keys, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		keys = append(keys, strconv.Itoa(a.PageSize))
	}

	return fmt.Sprintf("article:list:%s", strings.Join(keys, ":"))
}
