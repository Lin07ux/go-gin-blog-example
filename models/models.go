package models

import (
	"database/sql"
	"fmt"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/soft_delete"
	"log"
	"os"
	"time"
)

type Model struct {
	ID int `gorm:"primary_key" json:"id"`
	CreatedAt int `json:"created_at"`
	ModifiedAt int `json:"modified_at"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at"`
}

var db *gorm.DB

func Setup() {
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		sec.Key("USER").String(),
		sec.Key("PASSWORD").String(),
		sec.Key("HOST").String(),
		sec.Key("PORT").String(),
		sec.Key("NAME").String(),
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
		  log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标、前缀和日志包含的内容）
		  logger.Config{
			SlowThreshold:            0 * time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info,    // 日志级别
			IgnoreRecordNotFoundError: true,           // 忽略 ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,           // 启用彩色打印
		  },
		),
	})

	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}

	var sqlDB *sql.DB
	sqlDB, err = db.DB()

	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	_ = db.Callback().Update().Before("gorm:update").Register("timestamp:before_update", updateTimeStampForUpdateCallback)
}

// updateTimeStampForUpdateCallback will set `ModifiedAt` when updating
func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if _, ok := db.Get("gorm:update_column"); !ok {
		db.Statement.SetColumn("ModifiedAt", time.Now().Unix())
	}
}
