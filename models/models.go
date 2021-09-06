package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"log"
	"time"
)

type Model struct {
	ID int `gorm:"primary_key" json:"id"`
	CreatedAt int `json:"created_at"`
	ModifiedAt int `json:"modified_at"`
	DeletedAt int `json:"deleted_at"`
}

var db *gorm.DB

func init() {
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	db, err = gorm.Open(
		sec.Key("TYPE").String(),
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			sec.Key("USER").String(),
			sec.Key("PASSWORD").String(),
			sec.Key("HOST").String(),
			sec.Key("PORT").String(),
			sec.Key("NAME").String(),
		),
	)

	if err != nil {
		log.Println(err)
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
}

// updateTimeStampForCreateCallback will set `CreatedAt`, `ModifiedAt` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if ! scope.HasError() {
		nowTime := time.Now().Unix()

		if createTimeField, ok := scope.FieldByName("CreatedAt"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedAt"); ok {
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedAt` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("ModifiedAt", time.Now().Unix())
	}
}

// deleteCallback will set `DeletedAt` if delete column is exists or just delete from db when deleting
func deleteCallback(scope *gorm.Scope) {
	if ! scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedAt")

		// 存在软删除字段，且进行软删除；否则采用硬删除
		if ! scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// addExtraSpaceIfExist will add blank before str when str is not empty
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}

	return ""
}