package logger

// ZapLogConfig zap配置
type ZapLogConfig struct {
	// 日志文件路径
	LogDir string `mapstructure:"log_dir"`
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
