package plugin

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"gorm.io/gorm"
	"iter"
)

// NewPluginFunc 插件初始化函数,主程序通过该方法来新建插件实例
type NewPluginFunc = func() Plugin

// Engine ZeroBot Engine实例接口
// 保留原有注释
type Engine interface {
	// UsePreHandler 向该 Engine 添加新 PreHandler(Rule),
	// 会在 Rule 判断前触发，如果 preHandler
	// 没有通过，则 Rule, Matcher 不会触发
	//
	// 可用于分群组管理插件等
	UsePreHandler(rules ...zero.Rule)
	// UseMidHandler 向该 Engine 添加新 MidHandler(Rule),
	// 会在 Rule 判断后， Matcher 触发前触发，如果 midHandler
	// 没有通过，则 Matcher 不会触发
	//
	// 可用于速率限制等
	UseMidHandler(rules ...zero.Rule)
	// UsePostHandler 向该 Engine 添加新 PostHandler(Rule),
	// 会在 Matcher 触发后触发，如果 PostHandler 返回 false,
	// 则后续的 post handler 不会触发
	//
	// 可用于反并发等
	UsePostHandler(handler ...zero.Handler)
	// On 添加新的指定消息类型的匹配器
	On(typ string, rules ...zero.Rule) *zero.Matcher
	// OnMessage 消息触发器
	OnMessage(rules ...zero.Rule) *zero.Matcher
	// OnNotice 系统提示触发器
	OnNotice(rules ...zero.Rule) *zero.Matcher
	// OnRequest 请求消息触发器
	OnRequest(rules ...zero.Rule) *zero.Matcher
	// OnMetaEvent 元事件触发器
	OnMetaEvent(rules ...zero.Rule) *zero.Matcher
	// OnPrefix 前缀触发器
	OnPrefix(prefix string, rules ...zero.Rule) *zero.Matcher
	// OnSuffix 后缀触发器
	OnSuffix(suffix string, rules ...zero.Rule) *zero.Matcher
	// OnCommand 命令触发器
	OnCommand(commands string, rules ...zero.Rule) *zero.Matcher
	// OnRegex 正则触发器
	OnRegex(regexPattern string, rules ...zero.Rule) *zero.Matcher
	// OnKeyword 关键词触发器
	OnKeyword(keyword string, rules ...zero.Rule) *zero.Matcher
	// OnFullMatch 完全匹配触发器
	OnFullMatch(src string, rules ...zero.Rule) *zero.Matcher
	// OnFullMatchGroup 完全匹配触发器组
	OnFullMatchGroup(src []string, rules ...zero.Rule) *zero.Matcher
	// OnKeywordGroup 关键词触发器组
	OnKeywordGroup(keywords []string, rules ...zero.Rule) *zero.Matcher
	// OnCommandGroup 命令触发器组
	OnCommandGroup(commands []string, rules ...zero.Rule) *zero.Matcher
	// OnPrefixGroup 前缀触发器组
	OnPrefixGroup(prefix []string, rules ...zero.Rule) *zero.Matcher
	// OnSuffixGroup 后缀触发器组
	OnSuffixGroup(suffix []string, rules ...zero.Rule) *zero.Matcher
	// OnShell shell命令触发器
	OnShell(command string, model interface{}, rules ...zero.Rule) *zero.Matcher
}

// Plugin 所有插件需实现的接口
type Plugin interface {
	// OnInit 初始化插件(任意有关插件功能逻辑应放在此处进行，而不是在 NewPluginFunc),在Bot运行前调用
	OnInit(engine Engine, env Env) error
	// OnBoot engine准备就绪后调用
	OnBoot()
	// OnHelp 插件帮助回调
	OnHelp(ctx *zero.Ctx)
	// Name 插件名称，应具有唯一性
	Name() string
	// Version 插件版本,使用 go-SemVer 语义化版本格式
	//  example:
	//  func (p *myPlugin) Version() string {
	//		return "v1.0.0"
	//}
	Version() string
}

type ConfigProvider interface {
	// ConfigModel 获取配置文件结构，使用JSON Schema标签
	ConfigModel() any
}

// Env 插件运行环境
type Env interface {
	// Set 设置环境变量
	Set(key string, value any)
	// Get 获取环境变量
	Get(key string) any
	// FilePath 获取插件数据目录(不存在时会自动创建)
	FilePath() (string, error)
	// GetConf 从配置文件获取配置
	GetConf(conf any) error
	// GetDB 获取数据库连接
	GetDB() (*gorm.DB, error)
	// UseBot 获取并使用当前机器人实例
	UseBot(h zero.Handler)
	// Groups 获取启用的群
	Groups() Groups
	// SuperUser 获取SuperUser
	SuperUser() Users
	// Error 提交错误(由上层框架决定如何处理这个错误)
	Error(ctx *zero.Ctx, err error)
	// GetPlugin 获取对应的插件实例,可通过反射调用其他插件的方法
	GetPlugin(name string) (p Plugin, ok bool)
	// IsDisable 判断是否被禁用
	IsDisable() bool
}

// Groups 启用的群
type Groups interface {
	// IsContains 是否包含
	IsContains(groupId int64) bool
	// Rule 判断是否是开启的群
	Rule() zero.Rule
	// RangeGroup 遍历所有启用的群
	RangeGroup() iter.Seq[int64]
}

// Users 用户集合
type Users interface {
	// IsContains 是否包含
	IsContains(userId int64) bool
	// Rule 判断是否对应用户
	Rule() zero.Rule
	// RangeUser 遍历所有用户
	RangeUser() iter.Seq[int64]
}
