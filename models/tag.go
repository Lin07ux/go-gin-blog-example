package models

type Tag struct {
	Model

	Name string `json:"name"`
	State int `json:"state"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
}

func (Tag) TableName() string {
	return "tags"
}

// 查询文章标签列表
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

// 查询标签总量
func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

// 判断标签是否存在
func ExistTagByName(name string) bool {
	var tag Tag

	db.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

// 新增标签
func AddTag(name string, state int, createdBy string) bool {
	db.Create(&Tag {
		Name: name,
		State: state,
		CreatedBy: createdBy,
	})

	return true
}