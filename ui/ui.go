package ui

import (
	"github.com/invopop/jsonschema"
)

// 本文件定义一组「带 UI 提示」的配置类型。每个类型在 JSON Schema 里写入一个
// "k-ui" 扩展字段（外加少量参数字段），后台 assets/js/schema-form.js 读取
// "k-ui" 后用对应控件渲染，未识别的则回退到按 type 的默认渲染。
//
// 用法：把插件配置结构体里相应字段的类型换成这里的类型即可，例如：
//
//	type Conf struct {
//	    Welcome   TextArea `json:"welcome"   jsonschema:"title=欢迎语,description=支持多行"`
//	    BotToken  Secret   `json:"bot_token" jsonschema:"title=机器人 Token"`
//	    Template  Code     `json:"template"  jsonschema:"title=消息模板" jsonschema_extras:"k-lang=go-template,k-rows=10"`
//	    Theme     Color    `json:"theme"     jsonschema:"title=主题色"`
//	    Cooldown  Slider   `json:"cooldown"  jsonschema:"title=冷却秒数,minimum=0,maximum=60"`
//	    Admins    Tags     `json:"admins"    jsonschema:"title=管理员 QQ"`
//	    Keywords  StrTags  `json:"keywords"  jsonschema:"title=触发关键词"`
//	}
//
// 注意：minimum/maximum/multipleOf 等校验字段、以及 title/description，都来自字段
// 上的 jsonschema 标签，会作为属性级关键字与这里返回的 schema 合并，前端会一并读取。
// 也可以完全不用自定义类型，直接在任意字段上写 jsonschema_extras:"k-ui=textarea"
// 达到同样效果——自定义类型只是为了可复用、类型安全。

// TextArea 渲染为多行文本域，适合欢迎语、公告等需要换行的纯文本。
type TextArea string

func (TextArea) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:   "string",
		Extras: map[string]any{"k-ui": "textarea"},
	}
}

// Secret 渲染为密码框，并带「显示/隐藏」切换，适合 token、密钥、密码等敏感串。
// 值本身仍是普通字符串，仅前端默认遮蔽显示。
type Secret string

func (Secret) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:   "string",
		Extras: map[string]any{"k-ui": "secret"},
	}
}

// Code 渲染为等宽代码文本域，适合模板、正则、脚本片段等。
// 可配合字段标签 jsonschema_extras:"k-lang=yaml,k-rows=12" 标注语言与行数
// （k-lang 目前只作为 data-lang 暴露，前端未来可据此接入高亮）。
type Code string

func (Code) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:   "string",
		Extras: map[string]any{"k-ui": "code"},
	}
}

// Color 渲染为颜色选择器 + 十六进制输入，取值形如 #5b53c4。
type Color string

func (Color) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:    "string",
		Pattern: "^#[0-9a-fA-F]{6}$",
		Extras:  map[string]any{"k-ui": "color"},
	}
}

// Slider 渲染为滑块 + 数字框（两者联动），适合有明确上下限的整数，
// 例如冷却时间、并发数、阈值。范围请用字段标签提供：
//
//	jsonschema:"minimum=0,maximum=100,multipleOf=5"
//
// 若需要小数滑块，把下面的 Type 改成 "number" 并自定义一个浮点版本即可。
type Slider int

func (Slider) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:   "integer",
		Extras: map[string]any{"k-ui": "slider"},
	}
}

// Tags 渲染为标签输入（回车添加、点 × 删除），底层是「数字数组」，
// 适合 QQ 号、群号这类 ID 列表（如管理员、白名单）。
type Tags []int64

func (Tags) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:   "array",
		Items:  &jsonschema.Schema{Type: "integer"},
		Extras: map[string]any{"k-ui": "tags"},
	}
}

// StrTags 与 Tags 相同的标签控件，但底层是「字符串数组」，
// 适合关键词、别名、域名白名单等文本列表。
type StrTags []string

func (StrTags) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:   "array",
		Items:  &jsonschema.Schema{Type: "string"},
		Extras: map[string]any{"k-ui": "tags"},
	}
}
