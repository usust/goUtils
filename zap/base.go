package logger

const (
	// LOG_LEVEL_DEBUG 表示输出 debug 及以上级别日志.
	LOG_LEVEL_DEBUG = "debug"

	// LOG_LEVEL_INFO 表示输出 info 及以上级别日志.
	LOG_LEVEL_INFO = "info"

	// LOG_LEVEL_WARN 表示输出 warn 及以上级别日志.
	LOG_LEVEL_WARN = "warn"

	// LOG_LEVEL_ERROR 表示输出 error 及以上级别日志.
	LOG_LEVEL_ERROR = "error"

	// LOG_LEVEL_DPANIC 表示输出 dpanic 及以上级别日志.
	LOG_LEVEL_DPANIC = "dpanic"

	// LOG_LEVEL_PANIC 表示输出 panic 及以上级别日志.
	LOG_LEVEL_PANIC = "panic"

	// LOG_LEVEL_FATAL 表示输出 fatal 级别日志.
	LOG_LEVEL_FATAL = "fatal"
)

// ZapLogConfig 定义 zap 日志初始化所需的配置.
//
// 字段上的 mapstructure tag 用于从 viper 等配置库中反序列化配置.
// 如果调用方没有设置某些字段, InitZapCore 会为它们填充默认值.
type ZapLogConfig struct {
	// LogDir 是日志文件输出目录.
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

// ZapOption 用于以函数选项的方式覆盖 ZapLogConfig 默认值.
type ZapOption func(*ZapLogConfig)

// ZapWithLogDir 设置日志文件输出目录.
func ZapWithLogDir(dir string) ZapOption {
	return func(z *ZapLogConfig) { z.LogDir = dir }
}

// ZapWithLevel 设置日志最小输出级别.
func ZapWithLevel(level string) ZapOption {
	return func(z *ZapLogConfig) { z.LogLevel = level }
}

// ZapWithMaxSize 设置单个日志文件的最大大小, 单位为 MB.
func ZapWithMaxSize(max int) ZapOption {
	return func(z *ZapLogConfig) { z.MaxSize = max }
}

// ZapWithMaxBackups 设置最多保留的旧日志文件数量.
func ZapWithMaxBackups(max int) ZapOption {
	return func(z *ZapLogConfig) { z.MaxBackups = max }
}

// ZapWithMaxAge 设置旧日志文件最多保留天数.
func ZapWithMaxAge(days int) ZapOption {
	return func(z *ZapLogConfig) { z.MaxAge = days }
}

// ZapWithIsCompress 设置是否压缩轮转后的旧日志文件.
func ZapWithIsCompress(isCompress bool) ZapOption {
	return func(z *ZapLogConfig) { z.IsCompress = isCompress }
}
