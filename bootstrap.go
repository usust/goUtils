package go_utils

import (
	"fmt"
)

// Hooks 定义应用启动过程中可以被替换的初始化步骤.
//
// 这里使用函数类型而不是直接在 Bootstrap 中写死具体实现, 主要是为了让启动流程更容易测试
// 和扩展: 生产环境可以使用真实的配置加载与数据库初始化函数, 单元测试中可以替换为 mock 函数,
// 特殊运行环境中也可以注入定制的初始化逻辑.
type Hooks struct {
	// LoadConfig 用于加载应用运行所需的配置.
	//
	// 默认实现是 LoadConfigFileInto. 数据库、日志或其他组件的初始化通常依赖配置内容,
	// 因此它会在 InitDB 之前执行.
	LoadConfig func(dst any) error

	// LoadLog 用于根据已加载的配置初始化日志组件.
	//
	// 默认实现是 LoadLog. 如果调用方不需要初始化日志, 可以不传这个 hook.
	LoadLog func(cfg any) error

	// InitDB 用于初始化数据库连接.
	//
	// 如果调用方不需要真实数据库连接, 可以不传这个 hook; 如果需要初始化数据库,
	// 可以传入真实初始化函数, 或者在测试中传入用于记录调用顺序的 mock 函数.
	// InitDB func() error
}

// Bootstrap 使用默认初始化函数启动应用.
//
// 默认启动流程为:
//  1. 调用 LoadConfigFileInto 将配置文件加载到 dst 指向的配置结构体中.
//  2. 调用 LoadLog 根据配置初始化日志组件.
//  3. 如果提供了数据库初始化函数, 则继续初始化数据库连接.
//
// cfg 必须是非 nil 指针. 配置结构体由调用方定义, Bootstrap 不需要提前知道配置文件中
// 有哪些字段.
//
// 如果任一步骤失败, Bootstrap 会返回带有阶段信息的错误, 方便调用方判断失败阶段.
func Bootstrap(cfg any) error {
	return InitWith(cfg, Hooks{
		LoadConfig: LoadConfigFileInto,
		LoadLog:    LoadLog,
	})
}

// InitWith 按固定顺序执行应用启动流程.
//
// 参数 h 提供每个启动步骤的具体实现. 这种设计把“启动顺序”和“具体初始化逻辑”解耦,
// 既保留统一的启动入口, 又允许调用方按需替换其中某些步骤.
//
// 执行顺序固定为:
//  1. 校验必要 hook 是否存在.
//  2. 执行 LoadConfig, 将配置加载到 dst 中.
//  3. 如果 LoadLog 不为空, 执行 LoadLog.
//  4. 如果 InitDB 不为空, 执行 InitDB.
//
// 先加载配置再初始化数据库, 是因为数据库初始化通常需要读取配置中的连接地址、端口、
// 用户名、密码、数据库名等信息.
func InitWith(cfg any, h Hooks) error {
	// 启动 hook 是函数类型; 如果为 nil, 直接调用会触发 panic.
	// 因此这里先做显式校验, 把编程错误转换成可读的 error 返回给调用方.
	if h.LoadConfig == nil {
		return fmt.Errorf("bootstrap hook LoadConfig is nil")
	}

	// 第一步加载配置.
	// 使用 %w 包装原始错误, 调用方可以通过 errors.Is 或 errors.As 继续识别底层错误.
	if err := h.LoadConfig(cfg); err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}

	if h.LoadLog != nil {
		// 第二步初始化日志. 日志配置来自刚刚加载完成的 cfg.
		if err := h.LoadLog(cfg); err != nil {
			return fmt.Errorf("load log failed: %w", err)
		}
	}

	// 所有启动步骤都成功后返回 nil.
	return nil
}
