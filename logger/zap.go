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

// 对外暴露的日志对象
var (
	Logger        = zap.NewNop()
	SugaredLogger = Logger.Sugar()
	loggerMu      sync.Mutex
)

// InitZapCore
func InitZapCore(encoder *zapcore.EncoderConfig, option ...ZapOption) error {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	var encoderConfig zapcore.EncoderConfig

	// 设置日志输出格式，如果没有传入则使用默认的配置
	if encoder == nil {
		encoderConfig = defaultEncoderConfig()
	} else {
		encoderConfig = *encoder
	}

	// the only object of the zapConfig struct
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
	if zapConf.LogDir == "" {
		zapConf.LogDir = "./log"
	}
	if err := os.MkdirAll(zapConf.LogDir, 0o755); err != nil {
		return err
	}
	minLevel, err := zapcore.ParseLevel(strings.ToLower(zapConf.LogLevel))
	if err != nil {
		return err
	}

	// 配置JSON编码器
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 根据不同等级创建不同的日志文件
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

	// 设置日志级别过滤器
	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= minLevel && l >= zapcore.InfoLevel && l < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= minLevel && l >= zapcore.ErrorLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= minLevel && l == zapcore.DebugLevel
	})

	// 控制台输出全部日志
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		minLevel,
	)

	// 合并所有Core
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

	Logger = newLogger
	SugaredLogger = Logger.Sugar()
	zap.ReplaceGlobals(newLogger)
	if oldLogger != nil {
		_ = oldLogger.Sync()
	}
	if oldSugaredLogger != nil {
		_ = oldSugaredLogger.Sync()
	}

	return nil
}

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

// Sync 确保所有日志都写入磁盘
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
