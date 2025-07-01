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
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.AddConfigPath("../../.")
	viper.AddConfigPath("./config")

	// 初始化日志配置和创建日志对象
	if err := initLogger(); err != nil {
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

	if err := InitLogger(); err != nil {
		err = fmt.Errorf("日志初始化失败, %s", err)
		log.Fatal(err)
		return err
	} else {
		SugaredLogger.Info("日志初始化成功")

	}
	return nil
}
