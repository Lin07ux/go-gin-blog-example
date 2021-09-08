package cache

import (
	"fmt"
	"github.com/lin07ux/go-gin-example/models"
	"github.com/lin07ux/go-gin-example/pkg/gredis"
	"github.com/lin07ux/go-gin-example/pkg/logging"
	"time"
)

type Article struct {
	ID    int
	TagId int
	State int

	PageNum  int
	PageSize int
}

const detailCacheKey = "article:%d"
//const listCacheKey = "article:list:%s"

// GetDetail get article detail form cache
func (a *Article) GetDetail() *models.Article {
	key := fmt.Sprintf(detailCacheKey, a.ID)

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
	key := fmt.Sprintf(detailCacheKey, a.ID)

	if err := gredis.Set(key, article, 3600 * time.Second); err != nil {
		logging.Warn(err)
		return false
	}

	return true
}
