package utils

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// InitUtilsFromConfigFile 加载配置文件
func InitUtilsFromConfigFile() {
	// 读取配置信息
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 运行时和可执行程序同目录
	viper.AddConfigPath(".")

	// 初始化日志配置和创建日志对象
	if err := initLoggerByConfigFile(); err != nil {
		return
	}
}

// initLoggerByConfigFile 初始化zap日志配置
func initLoggerByConfigFile() error {
	if err := viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("error reading config file: %s", err.Error())
		log.Fatal(err)
		return err
	}

	if err := viper.Sub("zap").Unmarshal(&zapConf); err != nil {
		err = fmt.Errorf("an error occurred during the configuration of zap. err: %s", err.Error())
		log.Fatal(err)
		return err
	}

	if err := InitZapCore(); err != nil {
		err = fmt.Errorf("日志初始化失败, %s", err)
		log.Fatal(err)
		return err
	} else {
		SugaredLogger.Info("日志初始化成功")
	}
	return nil
}

// initGorm 初始化Gorm配置
//func initGorm() error {
//	if err := viper.ReadInConfig(); err != nil {
//		err = fmt.Errorf("error reading config file: %s", err.Error())
//		SugaredLogger.Errorln(err.Error())
//		return err
//	}
//
//	if err := viper.Sub("mysql").Unmarshal(&gormConf); err != nil {
//		err = fmt.Errorf("an error occurred during the configuration of zap. err: %s", err.Error())
//		SugaredLogger.Errorln(err.Error())
//		return err
//	}
//
//	if err := initGormConnect(); err != nil {
//		return err
//	} else {
//		SugaredLogger.Info("gorm初始化成功")
//	}
//	return nil
//}
