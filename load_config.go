package go_utils

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// ConfigValidator 是配置结构体可选实现的校验接口.
//
// 调用方传入的配置结构体如果实现了 Validate 方法, LoadConfigFileInto 会在反序列化成功后
// 自动调用该方法做业务校验; 如果没有实现, 则只负责读取和解析配置文件.
type ConfigValidator interface {
	Validate() error
}

// buildViper 根据配置文件路径创建一个独立的 viper 实例.
//
// 如果 configFile 不为空, 说明调用方已经通过环境变量指定了配置文件路径,
// 此时直接使用 SetConfigFile 精确读取该文件.
//
// 如果 configFile 为空, 则使用默认配置名和默认搜索路径查找 config.yaml.
// 这里每次创建新的 viper 实例, 避免使用全局 viper 带来的配置污染.
func buildViper(configFile string) *viper.Viper {
	v := viper.New()
	if configFile != "" {
		v.SetConfigFile(configFile)
		return v
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	return v
}

// LoadConfigFileInto 加载配置文件, 并反序列化到调用方传入的配置结构体中.
//
// 参数 cfg 必须是非 nil 指针, 例如:
//
//	var cfg AppConfig
//	if err := LoadConfigFileInto(&cfg); err != nil {
//	    return err
//	}
//
// 这个函数不要求 go-utils 预先知道配置文件中有哪些字段. 配置文件结构由调用方自己的
// struct tag 决定, viper 会按照传入结构体的字段定义进行映射.
//
// 加载规则:
//  1. 优先读取 ConfigFileEnv 指定的配置文件路径.
//  2. 如果环境变量为空, 则在默认搜索路径中查找 config.yaml.
//
// 如果 dst 实现了 ConfigValidator 接口, 配置解析成功后会自动调用 Validate 做字段校验.
func LoadConfigFileInto(cfg any) error {
	if cfg == nil {
		return fmt.Errorf("config destination is nil")
	}
	dstValue := reflect.ValueOf(cfg)
	if dstValue.Kind() != reflect.Ptr || dstValue.IsNil() {
		return fmt.Errorf("config destination must be a non-nil pointer")
	}

	// 允许环境变量前后带空白字符, 这里统一裁剪后再交给 viper.
	configFile := strings.TrimSpace(os.Getenv(ConfigFileEnv))
	fmt.Printf("config file env %s=%q\n", ConfigFileEnv, configFile)
	v := buildViper(configFile)

	// ReadInConfig 负责查找并读取配置文件.
	// 使用 %w 包装错误, 便于上层继续识别 viper 返回的具体错误类型.
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file failed (env %s=%q): %w", ConfigFileEnv, configFile, err)
	}
	fmt.Printf("load config file : %s\n", v.ConfigFileUsed())

	// 将配置文件内容映射到调用方传入的结构体, 后续代码就不需要直接依赖 viper 查询字段.
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshal config file failed: %w", err)
	}

	if validator, ok := cfg.(ConfigValidator); ok {
		// 尽早校验配置完整性, 避免程序启动到后续阶段才因为缺少字段失败.
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}
	}

	return nil
}
