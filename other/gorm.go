package other

import (
	"fmt"
	"go_utils/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormConfig struct {
	// mysql主机地址
	Path string `mapstructure:"path"`
	// mysql服务的端口号
	Port int `mapstructure:"port"`
	// mysql的用户名
	Username string `mapstructure:"username"`
	// 用户名对应的密码
	Password string `mapstructure:"password"`
	// 工具箱使用的数据库
	Database string `mapstructure:"database"`
	// 字符集，默认即可
	Charset string `mapstructure:"charset"`
	// 是否解析时间，默认true
	ParseTime string `mapstructure:"parse_time"`
	// 地区，默认即可
	Loc string `mapstructure:"loc"`
}

// gormConf gorm配置信息，通过viper读取config.yaml中内容
var gormConf GormConfig

// GormDB global gorm数据库对象
var GormDB *gorm.DB

// initGormConnect 初始化Gorm数据库连接
// @return error 抛出错误
func initGormConnect() error {
	// 数据库连接信息
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		gormConf.Username,
		gormConf.Password,
		gormConf.Path,
		gormConf.Port,
		gormConf.Database,
		gormConf.Charset,
		gormConf.ParseTime,
		gormConf.Loc)
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		logger.SugaredLogger.Errorw("gorm init err", "link", dsn, "err", err)
		return fmt.Errorf("gorm初始化时[%s]发生错误: %s", dsn, err)
	} else {
		// 成功连接到数据库
		GormDB = db
		logger.SugaredLogger.Info("gorm连接数据库成功")
	}

	// 数据迁移
	if err := GormDB.AutoMigrate(
	// 这里添加数据模型结构体，如：&models.User{},
	); err != nil {
		logger.SugaredLogger.Errorf("gorm auto migrate err: %v", err)
		return fmt.Errorf("gorm自动迁移时发生错误: %s", err.Error())
	} else {
		logger.SugaredLogger.Info("gorm自动迁移数据成功")
	}
	return nil
}
