package go_utils

import (
	"fmt"
	"reflect"

	myZaplog "github.com/usust/goUtils/zap"
)

// ZapLogConfig 是 zap 日志初始化配置.
//
// 这里使用类型别名对外暴露 zap 子包中的配置类型, 调用方只需要导入根包 go_utils,
// 不需要直接依赖内部的 zap 子包路径.
type ZapLogConfig = myZaplog.ZapLogConfig

// ZapLogConfigProvider 是应用配置结构体可选实现的日志配置接口.
//
// 如果调用方的配置结构体实现了这个接口, LoadLog 会优先通过该接口获取 zap 配置.
// 这比依赖固定字段名更清晰, 也更适合配置结构复杂的项目.
type ZapLogConfigProvider interface {
	ZapLogConfig() ZapLogConfig
}

// LoadLog 从应用配置中读取 zap 日志配置, 并注册日志组件.
//
// 这个函数用于放进 Hooks 中, 在配置文件加载完成后执行. cfg 通常就是调用
// LoadConfigFileInto 时传入的同一个配置结构体指针.
//
// LoadLog 获取日志配置的顺序:
//  1. 如果 cfg 实现了 ZapLogConfigProvider, 使用 cfg.ZapLogConfig() 的返回值.
//  2. 否则尝试读取 cfg.ZapLog 字段.
//  3. 否则尝试读取 cfg.ZapLogConfig 字段.
func LoadLog(cfg any) error {
	zapConf, err := resolveZapLogConfig(cfg)
	if err != nil {
		return err
	}
	return registerZap(zapConf)
}

// registerZap 使用配置结构体注册并初始化 zap 日志.
//
// 日志目录、日志级别、文件轮转策略等参数全部来自 conf,
// 避免调用方同时面对“配置文件”和“函数选项”两套配置入口.
//
// 调用成功后会同步更新根包暴露的 Logger、ZapLogger 和 SugaredLogger,
// 后续业务代码可以统一通过 go_utils.Logger 或 go_utils.SugaredLogger 写日志.
func registerZap(conf ZapLogConfig) error {
	if err := myZaplog.InitZapWithConfig(nil, conf); err != nil {
		return err
	}
	Logger = myZaplog.Logger
	ZapLogger = myZaplog.Logger
	SugaredLogger = myZaplog.SugaredLogger
	return nil
}

// SyncZap 刷新 zap 日志缓冲.
//
// 建议在程序退出前调用, 尽量确保日志写入磁盘或标准输出.
func SyncZap() {
	myZaplog.Sync()
}

func resolveZapLogConfig(cfg any) (ZapLogConfig, error) {
	if cfg == nil {
		return ZapLogConfig{}, fmt.Errorf("log config source is nil")
	}
	if provider, ok := cfg.(ZapLogConfigProvider); ok {
		return provider.ZapLogConfig(), nil
	}

	value := reflect.ValueOf(cfg)
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return ZapLogConfig{}, fmt.Errorf("log config source is nil pointer")
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return ZapLogConfig{}, fmt.Errorf("log config source must be a struct or implement ZapLogConfigProvider")
	}

	for _, fieldName := range []string{"ZapLog", "ZapLogConfig"} {
		field := value.FieldByName(fieldName)
		if !field.IsValid() {
			continue
		}
		if !field.CanInterface() {
			return ZapLogConfig{}, fmt.Errorf("log config field %s cannot be accessed", fieldName)
		}
		if fieldValue, ok := field.Interface().(ZapLogConfig); ok {
			return fieldValue, nil
		}
		if fieldValue, ok := field.Interface().(*ZapLogConfig); ok && fieldValue != nil {
			return *fieldValue, nil
		}
	}

	return ZapLogConfig{}, fmt.Errorf("log config not found: implement ZapLogConfigProvider or define ZapLog/ZapLogConfig field")
}
