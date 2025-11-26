# KohmeBot PluginV2

## 前言
本仓库为KohmeBot的插件定义，KohmeBot可通过`配置声明`来实现动态加载插件。 <br>

## ZeroBot
KohmeBot的底层是[ZeroBot](https://github.com/wdvxdr1123/ZeroBot.git)框架，提供了统一的插件管理和装配能力。<br>
在开发插件前，需要先熟悉`ZeroBot`框架的使用方法

## 插件仓库结构
```go
myplugin/ // 你的插件名称
├── go.mod // 模块入口
└── myplugin/ // 子包(和插件名称同名)
    └── plugin.go // 插件的实现
```
可查看[实例插件仓库](https://github.com/Kohmebot/chatai)。

## 生成插件仓库模板
```shell
# 安装kohme-gen
go install github.com/kohmebot/plugin/v2/cmd/kohme-gen@latest
# -n 指定你的插件名称 -r 指定模块名称,也就是github仓库地址
kohme-gen -n myplugin -r github.com/kohmebot/myplugin
```

- myplugin/plugin.go
```go
// myplugin/plugin.go
package myplugin
import "github.com/kohmebot/plugin/v2"

// 你的插件实现
type MyPluginImpl struct {
	// ...
}
// NewPlugin 初始化插件实例,方法签名必须固定
func NewPlugin() plugin.Plugin {
	return new(MyPluginImpl)
}
```


## 接口定义

### Plugin
插件接口，实现该接口，可以被KohmeBot插件系统正确加载
```go
type Plugin interface {...}
```

#### Init
`Init`方法会在Bot运行前调用，用于初始化插件，例如注册命令，事件等。
```go
// Init 初始化插件(任意有关插件功能逻辑应放在此处进行，而不是在 NewPluginFunc),在Bot运行前调用
Init(engine *plugin.Engine, env plugin.Env) error
```

#### OnBoot
`OnBoot`方法会在engine准备就绪后调用，注意不要阻塞
```go
// OnBoot engine准备就绪后调用
OnBoot()
```

#### OnHelp
`OnHelp`方法是用户在执行 `/help <plugin name>` 命令的回调，可用于发送插件帮助信息
```go
// OnHelp 插件帮助回调
OnHelp(ctx *zero.Ctx)
```

#### Name
`Name`方法用于返回插件名称，每个插件应具有唯一性
```go
// Name 插件名称，应具有唯一性
Name() string
```

#### Version
`Version`方法用于返回插件版本号，使用`go-SemVer`语义化版本格式<br>
```go
// Version 插件版本,使用 go-SemVer 语义化版本格式
//  example:
//  func (p *myPlugin) Version() string {
//		return "v1.0.0"
//}
Version() string
```


### Env
`Env`是插件的运行环境
```go
type Env interface {...}
```
#### Set
`Set` 方法用于设置插件运行环境变量，key-value键值对
```go
// Set 设置环境变量
Set(key string, value any)
```
#### Get
`Get`方法用于获取插件运行环境变量，通过key获取，取决于`kohmebot`的`plugins.yaml`的`plugins`配置，或者通过`Set`方法设置的变量
```go
// Get 获取环境变量
Get(key string) any
```
示例：
```yaml
# plugins.yaml
env:
  api_key: "your-api-key"
plugins:
  myplugin:
    conf:
      say: "hello world"
      time_duration: 10
```
```go
target,ok := env.Get("api_key").(string)

// do something...
```
#### FilePath
`FilePath` 获取插件的数据目录
目录路径是静态的，建议在`Init`方法中获取并保存
```go
// FilePath 获取插件数据目录(不存在时会自动创建)
FilePath() (string, error)
```

#### GetConf
`GetConf` 从`plugins.yaml`中对应插件的`conf`字段中解析相应的配置<br>
建议在`Init`方法中解析并保存
```go
// GetConf 从配置文件获取配置
GetConf(conf any) error
```
示例：
```yaml
# plugins.yaml
plugins:
  myplugin:
    target: 123456
    conf:
      say: "hello world"
      time_duration: 10
```
```go
// Config 结构体
type Config struct{
	Say string `yaml:"say"`
	TimeDuration int64 `yaml:"time_duration"`
}
//...
conf := Config{}
err := env.GetConf(&conf)
if err!=nil{return err}
```

#### GetDB
`GetDB` 获取插件的数据库连接(每个插件会有独立的连接池)，建议在`Init`方法中获取并保存
```go
GetDB() (*gorm.DB, error)
```
#### UseBot
`UseBot` 获取并使用当前bot实例
```go
// UseBot 获取并使用当前机器人实例
UseBot(h zero.Handler)
```

#### Groups
`Groups` 获取启用的群

#### SuperUsers
`SuperUser` 获取所有超级管理员

#### Error
`Error` 插件运行时抛出的错误
```go
// Error 提交错误(由上层框架决定如何处理这个错误)
Error(ctx *zero.Ctx, err error)
```

#### GetPlugin
`GetPlugin` 获取对应名称的插件实例<br>
若对方实例方法导出，则可以通过反射调用对应方法，提供插件间调用能力
```go
// GetPlugin 获取对应的插件实例,可通过反射调用其他插件的方法
GetPlugin(name string) (p Plugin, ok bool)
```

#### IsDisable
`IsDisable`判断插件功能此时是否被禁用

### Groups
`Groups`是一个已启用群的集合
```go
type Groups interface {...}
```

#### IsContains
`IsContains` 已启用群集合内是否包含该群，也就是说，判断群是否启用

#### Rule
`Rule`在engine.On时可注入的判断规则，将判断是否是群消息，且群是否启用<br>

#### RangeGroup
`RangeGroup` 遍历所有已启用的群

### Users
`Users`是一个用户的集合，方法使用与`Groups`相同
```go
type Users interface {...}
```

## 自动构建
1. 将你的代码推送到仓库,并打包对应版本tag
2. 在主程序的`plugins.yaml`中,指定插件名称与插件仓库地址,使用构建脚本即可自动下载插件

