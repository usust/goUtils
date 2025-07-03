package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// InitConfig 加载配置文件
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// 当前路径为项目执行路径，默认为当前项目路径
	viper.AddConfigPath("./pkg/goUtils")
	viper.AddConfigPath(".")

	// 初始化日志配置和创建日志对象
	if err := initLogger(); err != nil {
		return
	}

	// 初始化gorm，完成数据迁移
	if err := initGorm(); err != nil {
		return
	}
}

// initLogger 初始化zap日志配置
func initLogger() error {
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

	if err := initZapCore(); err != nil {
		err = fmt.Errorf("日志初始化失败, %s", err)
		log.Fatal(err)
		return err
	} else {
		SugaredLogger.Info("日志初始化成功")
	}
	return nil
}

// initGorm 初始化Gorm配置
func initGorm() error {
	if err := viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("error reading config file: %s", err.Error())
		SugaredLogger.Errorln(err.Error())
		return err
	}

	if err := viper.Sub("mysql").Unmarshal(&gormConf); err != nil {
		err = fmt.Errorf("an error occurred during the configuration of zap. err: %s", err.Error())
		SugaredLogger.Errorln(err.Error())
		return err
	}

	if err := initGormConnect(); err != nil {
		return err
	} else {
		SugaredLogger.Info("gorm初始化成功")
	}
	return nil
}
