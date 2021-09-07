package models

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
func ExistArticleById(id int) bool {
	var article Article

	db.Select("id").Where("id = ?", id).First(&article)

	return article.ID > 0
}

// 获取文章的数量
func GetArticleTotal(maps interface{}) (count int64) {
	db.Model(&Article{}).Where(maps).Count(&count)

	return
}

// 获取文章列表
func GetArticles(pageNum, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

// 获取文章内容
func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	_ = db.Model(&article).Association("Tag").Find(&article.Tag)

	return
}

// 添加文章
func AddArticle(article *Article) bool {
	db.Create(article)

	return true
}

// 编辑文章
func EditArticle(id int, article *Article) bool {
	data := make(map[string]interface{})
	data["modified_by"] = article.ModifiedBy

	if article.TagID > 0 {
		data["tag_id"] = article.TagID
	}
	if article.Title != "" {
		data["title"] = article.Title
	}
	if article.Cover != "" {
		data["cover"] = article.Cover
	}
	if article.Desc != "" {
		data["desc"] = article.Desc
	}
	if article.Content != "" {
		data["content"] = article.Content
	}
	if article.State >= 0 {
		data["state"] = article.State
	}

	db.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

// 删除文章
func DeleteArticle(id int) bool {
	db.Delete(&Article{}, id)

	return true
}

// 清理全部的已删除文章
func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_at > ?", 0).Delete(&Article{})

	return true
}
