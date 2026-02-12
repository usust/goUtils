package logger

const (
	LOG_LEVEL_DEBUG  = "debug"
	LOG_LEVEL_INFO   = "info"
	LOG_LEVEL_WARN   = "warn"
	LOG_LEVEL_ERROR  = "error"
	LOG_LEVEL_DPANIC = "dpanic"
	LOG_LEVEL_PANIC  = "panic"
	LOG_LEVEL_FATAL  = "fatal"
)

// ZapLogConfig custom zap config
type ZapLogConfig struct {
	// the dir of log
	LogDir string `mapstructure:"log_dir"`
	// 日志最小输出级别
	LogLevel string `mapstructure:"log_level"`
	// 单个日志文件的最大大小 (MB)
	MaxSize int `mapstructure:"max_size"`
	// 最多保留的旧日志文件个数
	MaxBackups int `mapstructure:"max_backups"`
	// 最多保留旧日志文件的天数
	MaxAge int `mapstructure:"max_age"`
	// 是否压缩旧日志文件
	IsCompress bool `mapstructure:"iscompress"`
}

type ZapOption func(*ZapLogConfig)

func ZapWithLogDir(dir string) ZapOption {
	return func(z *ZapLogConfig) { z.LogDir = dir }
}

func ZapWithLevel(level string) ZapOption {
	return func(z *ZapLogConfig) { z.LogLevel = level }
}

func ZapWithMaxSize(max int) ZapOption {
	return func(z *ZapLogConfig) { z.MaxSize = max }
}

func ZapWithMaxBackups(max int) ZapOption {
	return func(z *ZapLogConfig) { z.MaxBackups = max }
}

func ZapWithMaxAge(days int) ZapOption {
	return func(z *ZapLogConfig) { z.MaxAge = days }
}

func ZapWithIsCompress(isCompress bool) ZapOption {
	return func(z *ZapLogConfig) { z.IsCompress = isCompress }
}
