package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"log"
)

var db *gorm.DB

type Model struct {
	ID int `gorm:"primary_key" json:"id"`
	CreateAt int `json:"created_at"`
	ModifiedAt int `json:"modified_at"`
}

func init() {
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	db, err = gorm.Open(
		sec.Key("TYPE").String(),
		fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			sec.Key("USER").String(),
			sec.Key("PASSWORD").String(),
			sec.Key("HOST").String(),
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
}

func CloseDB() {
	defer db.Close()
}