package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type app struct {
	JwtSecret string
	PageSize int
	RuntimeRootPath string

	LogSavePath string
	LogSaveName string
	LogFileExt string
	TimeFormat string
}

type server struct {
	RunMode string
	HttpPort int
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

type database struct {
	Type string
	Host string
	Port int
	Name string
	User string
	Password string
}

// 应用配置
var AppSetting = &app{}
// HTTP 服务器配置
var ServerSetting = &server{}
// 数据库配置
var DatabaseSetting = &database{}

// Setup initialize settings from config file
func Setup() {
	config, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	err = config.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("config.MapTo AppSetting err: %v", err)
	}

	err = config.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("config.MapTo ServerSetting err: %v", err)
	}

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	err = config.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("config.MapTo DatabaseSetting err: %v", err)
	}
}
