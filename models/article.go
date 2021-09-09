package models

import "gorm.io/gorm"

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title string `json:"title"`
	Cover string `json:"cover"`
	Desc string `json:"desc"`
	Content string `json:"content"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

// 使用 ID 判断文章是否存在
func ExistArticleById(id int) (bool, error) {
	var article Article

	err := db.Select("id").Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

// 获取文章的数量
func GetArticleTotal(maps interface{}) (count int64, err error) {
	if err = db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// 获取文章列表
func GetArticles(offset, limit int, maps interface{}) (articles []*Article, err error) {
	err = db.Preload("Tag").Where(maps).Offset(offset).Limit(limit).Find(&articles).Error

	return articles, err
}

// 获取文章内容
func GetArticle(id int) (*Article, error) {
	article := &Article{}

	err := db.Where("id = ?", id).First(article).Error
	if err == nil {
		err = db.Model(&article).Association("Tag").Find(&article.Tag)
		if err == nil {
			return article, nil
		}
	} else if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return nil, err
}

// 添加文章
func AddArticle(article *Article) (int, error) {
	if err := db.Create(article).Error; err != nil {
		return 0, err
	}

	return article.ID, nil
}

// 编辑文章
func EditArticle(id int, data interface{}) error {
	return db.Model(&Article{}).Where("id = ?", id).Updates(data).Error
}

// 删除文章
func DeleteArticle(id int) error {
	return db.Delete(&Article{}, id).Error
}

// 清理全部的已删除文章
func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_at > ?", 0).Delete(&Article{})

	return true
}
