package go_utils

import "go.uber.org/zap"

var (
	// ConfigFileEnv 是用于指定配置文件路径的环境变量名.
	//
	// 调用方可以通过设置 EXPORT_CONFIG_FILE 来覆盖默认配置文件位置,
	// 让程序在不同运行环境中加载不同的配置文件.
	ConfigFileEnv = "EXPORT_CONFIG_FILE"
	// Logger 是 go-utils 对外暴露的 zap logger.
	//
	// 初始值为 zap.NewNop(), 在调用 LoadLog 或 RegisterZap 成功后会替换为真实 logger.
	Logger = zap.NewNop()

	// ZapLogger 是 Logger 的兼容别名.
	ZapLogger = Logger

	// SugaredLogger 是 Logger 的 SugaredLogger 版本, 适合使用格式化或键值对日志.
	SugaredLogger = Logger.Sugar()
)
