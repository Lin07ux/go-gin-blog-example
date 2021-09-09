package cache

import (
	"fmt"
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/gredis"
	"github.com/lin07ux/go-gin-example/pkg/logging"
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

const articleDetailCacheKey = "article:%d"
const articleListCacheKey = "article:list:%s"

// GetDetail get article detail form cache
func (a *Article) GetDetail() *models.Article {
	key := fmt.Sprintf(articleDetailCacheKey, a.ID)

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
	key := fmt.Sprintf(articleDetailCacheKey, a.ID)

	if err := gredis.Set(key, article, 3600 * time.Second); err != nil {
		logging.Warn(err)
		return false
	}

	return true
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

	if err := gredis.Set(key, articles, 3600 * time.Second); err != nil {
		logging.Warn(err)
		return false
	}

	return true
}

// getArticleListKey build the articles list cache key by query form
func (a *Article) getArticleListKey() string {
	var keys []string

	if a.ID > 0 {
		keys = append(keys, strconv.Itoa(a.ID))
	}
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

	return fmt.Sprintf(articleListCacheKey, strings.Join(keys, ":"))
}
