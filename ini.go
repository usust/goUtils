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

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %s", err)
		return
	}

	if err := viper.Sub("zap").Unmarshal(&zapConf); err != nil {
		return log.Fatalf("an error occurred during the configuration of zap. err: %s", err.Error())
	}

	if err := viper.Unmarshal(&GConfig); err != nil {
		log.Fatalf("Error unmarshalling config, %s", err)
		return
	}

	// 初始化日志
	if err := InitLogger(&GConfig.LogCfg); err != nil {
		log.Fatalf("日志初始化失败, %s", err)
		return
	} else {
		SugaredLogger.Info("日志初始化成功")
	}

	// 初始化GORM
	if err := InitGorm(); err != nil {
		SugaredLogger.Errorw("初始化数据存储失败", "err", err)
	} else {
		SugaredLogger.Info("数据存储初始化完成")
	}
}
