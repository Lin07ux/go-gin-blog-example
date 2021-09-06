package models

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title string `json:"title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

// 数据表名
func (Article) TableName() string {
	return "articles"
}

// 使用 ID 判断文章是否存在
func ExistArticleById(id int) bool {
	var article Article

	db.Select("id").Where("id = ?", id).First(&article)

	return article.ID > 0
}

// 获取文章的数量
func GetArticleTotal(maps interface{}) (count int) {
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
	db.Model(&article).Related(&article.Tag)

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
	if article.Desc != "" {
		data["desc"] = article.Desc
	}
	if article.Content != "" {
		data["content"] = article.Content
	}
	if article.State >= 0 {
		data["state"] = article.State
	}

	db.Model(&Article{}).Where("id = ?", id).Update(data)

	return true
}

// 删除文章
func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})

	return true
}
