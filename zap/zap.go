package logger

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 是当前包对外暴露的结构化日志对象.
//
// 初始值是 zap.NewNop(), 即不会输出任何日志. 调用 InitZapCore 成功后,
// Logger 会被替换为真实的 zap.Logger.
var (
	Logger        = zap.NewNop()
	SugaredLogger = Logger.Sugar()
	loggerMu      sync.Mutex
)

// InitZapCore 初始化 zap 日志核心.
//
// encoder 用于定制日志字段名、时间格式、级别格式等编码配置; 如果传入 nil,
// 则使用 defaultEncoderConfig 返回的默认配置.
//
// option 用于覆盖日志目录、最小日志级别、日志轮转策略等配置. 未指定的配置会使用默认值:
//   - 日志目录: ./log
//   - 最小日志级别: debug
//   - 单文件最大大小: 10 MB
//   - 最多保留旧文件数量: 30
//   - 最多保留旧文件天数: 180
//   - 轮转后压缩旧文件: true
//
// 初始化成功后, 会同时更新本包的 Logger / SugaredLogger, 并通过 zap.ReplaceGlobals
// 替换 zap.L() 和 zap.S() 使用的全局 logger.
func InitZapCore(encoder *zapcore.EncoderConfig, option ...ZapOption) error {
	// 先构造默认配置, 再按 option 覆盖调用方指定的字段.
	zapConf := ZapLogConfig{
		LogDir:     "./log",
		LogLevel:   LOG_LEVEL_DEBUG,
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     180,
		IsCompress: true,
	}
	for _, opt := range option {
		if opt != nil {
			opt(&zapConf)
		}
	}
	return InitZapWithConfig(encoder, zapConf)
}

// InitZapWithConfig 使用完整配置初始化 zap 日志核心.
//
// 该函数适合配置来自配置文件的场景. 调用方先把配置文件反序列化为 ZapLogConfig,
// 再将该配置传入本函数, 避免同时使用函数选项和配置文件两套入口.
func InitZapWithConfig(encoder *zapcore.EncoderConfig, zapConf ZapLogConfig) error {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	var encoderConfig zapcore.EncoderConfig

	// 设置日志输出格式; 如果没有传入则使用默认配置.
	if encoder == nil {
		encoderConfig = defaultEncoderConfig()
	} else {
		encoderConfig = *encoder
	}

	if zapConf.LogDir == "" {
		zapConf.LogDir = "./log"
	}
	if zapConf.LogLevel == "" {
		zapConf.LogLevel = LOG_LEVEL_DEBUG
	}
	if zapConf.MaxSize <= 0 {
		zapConf.MaxSize = 10
	}
	if zapConf.MaxBackups <= 0 {
		zapConf.MaxBackups = 30
	}
	if zapConf.MaxAge <= 0 {
		zapConf.MaxAge = 180
	}
	// 确保日志目录存在; lumberjack 不负责创建父目录.
	if err := os.MkdirAll(zapConf.LogDir, 0o755); err != nil {
		return err
	}
	// 将字符串日志级别转换成 zapcore.Level, 无效级别会直接返回错误.
	minLevel, err := zapcore.ParseLevel(strings.ToLower(zapConf.LogLevel))
	if err != nil {
		return err
	}

	// 文件日志使用 JSON 编码, 便于后续被日志平台采集和解析.
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 根据不同等级创建不同的日志文件, 并交给 lumberjack 处理文件轮转.
	infoWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(zapConf.LogDir, "info.log"),
		MaxSize:    zapConf.MaxSize,
		MaxBackups: zapConf.MaxBackups,
		MaxAge:     zapConf.MaxAge,
		Compress:   zapConf.IsCompress,
	})
	errorWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(zapConf.LogDir, "error.log"),
		MaxSize:    zapConf.MaxSize,
		MaxBackups: zapConf.MaxBackups,
		MaxAge:     zapConf.MaxAge,
		Compress:   zapConf.IsCompress,
	})
	debugWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(zapConf.LogDir, "debug.log"),
		MaxSize:    zapConf.MaxSize,
		MaxBackups: zapConf.MaxBackups,
		MaxAge:     zapConf.MaxAge,
		Compress:   zapConf.IsCompress,
	})

	// 设置日志级别过滤器:
	//   - info.log 记录 info/warn 等非 error 级别日志.
	//   - error.log 记录 error 及以上级别日志.
	//   - debug.log 只记录 debug 日志.
	// 三者都会受到 minLevel 的统一控制.
	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= minLevel && l >= zapcore.InfoLevel && l < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= minLevel && l >= zapcore.ErrorLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= minLevel && l == zapcore.DebugLevel
	})

	// 控制台输出所有达到 minLevel 的日志, 方便本地开发和容器标准输出采集.
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		minLevel,
	)

	// 合并多个 Core, 让同一条日志按级别同时写入对应文件和控制台.
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, infoWriter, infoLevel),
		zapcore.NewCore(jsonEncoder, errorWriter, errorLevel),
		zapcore.NewCore(jsonEncoder, debugWriter, debugLevel),
		consoleCore,
	)

	newLogger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	oldLogger := Logger
	oldSugaredLogger := SugaredLogger

	// 替换本包日志对象和 zap 全局日志对象.
	Logger = newLogger
	SugaredLogger = Logger.Sugar()
	zap.ReplaceGlobals(newLogger)
	// 尝试刷新旧 logger 中尚未写出的日志, 忽略 Sync 错误以免影响本次初始化结果.
	if oldLogger != nil {
		_ = oldLogger.Sync()
	}
	if oldSugaredLogger != nil {
		_ = oldSugaredLogger.Sync()
	}

	return nil
}

// defaultEncoderConfig 返回默认日志编码配置.
//
// 默认配置会输出时间、级别、logger 名称、调用位置、消息和 stacktrace.
// 时间格式使用本地易读的 "2006-01-02 15:04:05".
func defaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// Sync 确保内存缓冲中的日志尽量刷新到底层 writer.
//
// 程序退出前可以调用该函数减少日志丢失概率. 某些平台上 stdout/stderr 的 Sync
// 可能返回无害错误, 当前实现会忽略这些错误.
func Sync() {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	if Logger == nil || SugaredLogger == nil {
		return
	}
	err := Logger.Sync()
	if err != nil {
		return
	}
	err = SugaredLogger.Sync()
	if err != nil {
		return
	}
}
