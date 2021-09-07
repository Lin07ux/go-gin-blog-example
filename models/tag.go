package models

type Tag struct {
	Model

	Name string `json:"name"`
	State int `json:"state"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
}

// 查询文章标签列表
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

// 查询标签总量
func GetTagTotal(maps interface{}) (count int64) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

// 判断标签名称是否存在
func ExistTagByName(name string) bool {
	var tag Tag

	db.Select("id").Where("name = ?", name).First(&tag)

	return tag.ID > 0
}

// 判断标签是否存在
func ExistTagById(id int) bool {
	var tag Tag

	db.Select("id").Where("id = ?", id).First(&tag)
	
	return tag.ID > 0
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

// 编辑标签
func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id = ?", id).Updates(data)

	return true
}

// 删除标签
func DeleteTag(id int) bool {
	db.Delete(&Tag{}, id)

	return true
}

// 清理全部的已删除标签
func CleanAllTag() bool {
	db.Unscoped().Where("deleted_at > ?", 0).Delete(&Tag{})

	return true
}