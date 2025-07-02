# goUtils

## 依赖项

### viper
``` go
go get github.com/spf13/viper 
```

解析配置文件和支持动态解析配置文件。

## 功能

1. 日志功能

  通过 `SugaredLogger` 和 `Logger` 提供快捷日志功能。默认情况下输出分别输出 `info.log`、`error.log` 和 `debug.log`。

  使用方法

```go
Logger.Info("日志初始化成功")
SugaredLogger.Info("日志初始化成功")
```
