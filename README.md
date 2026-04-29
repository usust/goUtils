# goUtils

## 依赖项

### viper
  ``` sh
    go get github.com/spf13/viper 
  ```

解析配置文件和支持动态解析配置文件。

### gorm

  ```sh
    go get "gorm.io/gorm"
    # 默认使用 MYSQL
    go get "gorm.io/driver/mysql"
  ```
  直接映射管理数据库

## 功能

### 日志功能

  通过 `SugaredLogger` 和 `Logger` 提供快捷日志功能。默认情况下输出分别输出 `info.log`、`error.log` 和 `debug.log`。

  使用方法
    
  ```go
    Logger.Info("日志初始化成功")
    SugaredLogger.Info("日志初始化成功")
  ```

### 数据库

  创建、管理数据库

  使用方法
  ``` go
    if err := GormDB.Create(asset).Error; err != nil {
		config.SugaredLogger.Errorw("创建数据资产记录失败", "asset", asset, "error", err)
		return errors.New("创建数据资产记录失败")
	}
	return nil
  ```

