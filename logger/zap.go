package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//type zapConfig struct {
//	// path of the log file
//	LogDir string `mapstructure:"log_dir"`
//	// max size (MB) of the single log file
//	MaxSize int `mapstructure:"max_size"`
//	// max number of saved log files
//	MaxBackups int `mapstructure:"max_backups"`
//	// max day number of the saved log file 最多保留旧日志文件的天数
//	MaxAge int `mapstructure:"max_age"`
//	// whether compress the log file
//	IsCompress bool `mapstructure:"compress"`
//}

// 对外暴露的日志对象
var (
	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
)

// InitZapCore 初始化zap日志
func InitZapCore(encoder *zapcore.EncoderConfig, option ...ZapOption) error {
	var encoderConfig zapcore.EncoderConfig

	// 设置日志输出格式
	if encoder == nil {
		encoderConfig = zapcore.EncoderConfig{
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
	} else {
		encoderConfig = *encoder
	}

	// 配置JSON编码器
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// the only object of the zapConfig struct
	zapConf := ZapLogConfig{
		LogDir:     "./log",
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     180,
		IsCompress: true,
	}

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
		return l >= zapcore.InfoLevel && l < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.ErrorLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == zapcore.DebugLevel
	})

	// 控制台输出全部日志
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	// 合并所有Core
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, infoWriter, infoLevel),
		zapcore.NewCore(jsonEncoder, errorWriter, errorLevel),
		zapcore.NewCore(jsonEncoder, debugWriter, debugLevel),
		consoleCore,
	)

	Logger = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	SugaredLogger = Logger.Sugar()
	zap.ReplaceGlobals(Logger)

	return nil
}

// Sync 确保所有日志都写入磁盘
func Sync() {
	err := Logger.Sync()
	if err != nil {
		return
	}
	err = SugaredLogger.Sync()
	if err != nil {
		return
	}
}
