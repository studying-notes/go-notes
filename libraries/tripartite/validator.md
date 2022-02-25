---
date: 2020-08-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "validator - 参数校验"  # 文章标题
url:  "posts/go/libraries/tripartite/validator"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go", "validator" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

- [官网](#官网)
- [安装](#安装)
- [标签详解](#标签详解)
	- [限制字符串长度](#限制字符串长度)
	- [定长字符串](#定长字符串)
	- [数值比较](#数值比较)
	- [列表长度](#列表长度)
	- [嵌套结构验证](#嵌套结构验证)
	- [必填项/条件可填](#必填项条件可填)
	- [自定义枚举值](#自定义枚举值)
	- [字段间关系](#字段间关系)
- [Gin 增加的标签](#gin-增加的标签)
	- [Content-Type](#content-type)
	- [字段的默认值](#字段的默认值)
- [全部标签](#全部标签)
	- [Fields](#fields)
	- [Network:](#network)
	- [Strings](#strings)
	- [Format](#format)
	- [Comparisons](#comparisons)
	- [Other](#other)
		- [Aliases](#aliases)

## 官网

```shell
https://github.com/go-playground/validator
```

## 安装

```shell
go get github.com/go-playground/validator/v10
```

```go
import "github.com/go-playground/validator/v10"
```

Gin 框架默认设置 `TagName` 为 `binding`。

```go
func Validator() *validator.Validate {
	vOnce.Do(func() {
		v = validator.New()
		v.SetTagName("binding")
	})

	return v
}
```

## 标签详解

### 限制字符串长度

2 <= length < 100

```go
`binding:"required,min=2,max=100" example:"限制字符串长度"`       
```

min 不能大于等于 max，否则将造成异常。

### 定长字符串

```go
`binding:"required,len=100" example:"定长字符串"`       
```

翻译中文比较特殊，len 必须大于 1，否则将造成异常。

### 数值比较

- `eq`: equal，等于；
- `ne`: not equal，不等于；
- `gt`: great than，大于；
- `gte`: great than equal，大于等于；
- `lt`: less than，小于；
- `lte`: less than equal，小于等于；

```go
// Pager 分页处理
type Pager struct {
	Page      int   `json:"page" form:"page" url:"page" binding:"gte=1"`                         // 页码
	PageSize  int   `json:"page_size" form:"page_size" url:"page_size" binding:"gte=10,lte=100"` // 每页数量
	Order     int   `json:"order" form:"order" binding:"oneof=0 1 2 3"`                          // 排序顺序
	TotalRows int64 `json:"total_rows"`                                                          // 总行数
}
```

### 列表长度

```go
type Pager struct {
	Array []string `binding:"required,gt=0"`
}      
```

### 嵌套结构验证

`dive` 一般用在 slice、array、map、嵌套的 struct 验证中，作为分隔符表示进入里面一层的验证规则。

```go
type Test struct {
	Array []string `validate:"required,gt=0,dive,max=2"`
    // gt=0 代表 field.len()>0

	Map map[string]string `validate:"required,gt=0,dive,keys,max=10,endkeys,required,max=2"`
	// dive 表示进入里面一层的验证，例如 a={"b":"c"} 中 dive 之前的 required 表示 a 是必填项，大于0，
	// dive 之后出现 keys 与 endkeys 之间代表验证 map 的 keys 的 tag 值：max=10，即长度不大于 10
	// 后面的是验证 value，required 必填并且最大长度是 2
}
```

### 必填项/条件可填

`required` 表示字段必须有值，并且不为默认值，例如 `bool` 默认值为 `false`、`string` 默认值为 `””`、`int` 默认值为 `0`。如果有些字段是可填的，并且需要满足某些规则的，那么需要使用 `omitempty`。

```go
type Test struct {
	Id         string `form:"charger_id" validate:"omitempty,uuid4"`
    // omitempty 表示变量可以不填，但是填的时候必须满足条件
	Password   string `form:"password" validate:"omitempty,min=5,max=128"`
}
```

### 自定义枚举值

oneof 自定义枚举值。

```go
type Test struct {
	Gender uint8  `json:"gender" binding:"oneof=0 1 2"`
}
```

### 字段间关系

- `eqfield=Field`: 必须等于 Field 的值；
- `nefield=Field`: 必须不等于 Field 的值；
- `gtfield=Field`: 必须大于 Field 的值；
- `gtefield=Field`: 必须大于等于 Field 的值；
- `ltfield=Field`: 必须小于 Field 的值；
- `ltefield=Field`: 必须小于等于 Field 的值；
- `eqcsfield=Other.Field`: 必须等于 struct Other 中 Field 的值；
- `necsfield=Other.Field`: 必须不等于 struct Other 中 Field 的值；
- `gtcsfield=Other.Field`: 必须大于 struct Other 中 Field 的值；
- `gtecsfield=Other.Field`: 必须大于等于 struct Other 中 Field 的值；
- `ltcsfield=Other.Field`: 必须小于 struct Other 中 Field 的值；
- `ltecsfield=Other.Field`: 必须小于等于 struct Other 中 Field 的值；
- `required_with=Field1 Field2`: 在 Field1 或者 Field2 存在时，必须；
- `required_with_all=Field1 Field2`: 在 Field1 与 Field2 都存在时，必须；
- `required_without=Field1 Field2`: 在 Field1 或者 Field2 不存在时，必须；
- `required_without_all=Field1 Field2`: 在 Field1 与 Field2 都不存在时，必须；

## Gin 增加的标签

### Content-Type

```go
var (
	JSON          = jsonBinding{} `json`
	XML           = xmlBinding{}`xml`
	Form          = formBinding{}`form`
	Query         = queryBinding{} `form`
	FormPost      = formPostBinding{}`form`
	FormMultipart = formMultipartBinding{}`form`
	ProtoBuf      = protobufBinding{}
	MsgPack       = msgpackBinding{}
)
```

常用的就 form、query、json 和 xml；经过测试 url 中的传参（application/x-www-form-urlencoded）和 body 中传参**都可以使用 form 标签**

```go
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}
```

### 字段的默认值

```go
type Login struct {
	User     string `form:"user,default=admin"  binding:"required"`
	Password string `form:"password"  binding:"required"`
}
```

## 全部标签

### Fields

| Tag | Description |
| - | - |
| eqcsfield | Field Equals Another Field (relative)|
| eqfield | Field Equals Another Field |
| fieldcontains | NOT DOCUMENTED IN doc.go |
| fieldexcludes | NOT DOCUMENTED IN doc.go |
| gtcsfield | Field Greater Than Another Relative Field |
| gtecsfield | Field Greater Than or Equal To Another Relative Field |
| gtefield | Field Greater Than or Equal To Another Field |
| gtfield | Field Greater Than Another Field |
| ltcsfield | Less Than Another Relative Field |
| ltecsfield | Less Than or Equal To Another Relative Field |
| ltefield | Less Than or Equal To Another Field |
| ltfield | Less Than Another Field |
| necsfield | Field Does Not Equal Another Field (relative) |
| nefield | Field Does Not Equal Another Field |

### Network:

| Tag | Description |
| - | - |
| cidr | Classless Inter-Domain Routing CIDR |
| cidrv4 | Classless Inter-Domain Routing CIDRv4 |
| cidrv6 | Classless Inter-Domain Routing CIDRv6 |
| datauri | Data URL |
| fqdn | Full Qualified Domain Name (FQDN) |
| hostname | Hostname RFC 952 |
| hostname_port | HostPort |
| hostname_rfc1123 | Hostname RFC 1123 |
| ip | Internet Protocol Address IP |
| ip4_addr | Internet Protocol Address IPv4 |
| ip6_addr | Internet Protocol Address IPv6 |
| ip_addr | Internet Protocol Address IP |
| ipv4 | Internet Protocol Address IPv4 |
| ipv6 | Internet Protocol Address IPv6 |
| mac | Media Access Control Address MAC |
| tcp4_addr | Transmission Control Protocol Address TCPv4 |
| tcp6_addr | Transmission Control Protocol Address TCPv6 |
| tcp_addr | Transmission Control Protocol Address TCP |
| udp4_addr | User Datagram Protocol Address UDPv4 |
| udp6_addr | User Datagram Protocol Address UDPv6 |
| udp_addr | User Datagram Protocol Address UDP |
| unix_addr | Unix domain socket end point Address |
| uri | URI String |
| url | URL String |
| url_encoded | URL Encoded |
| urn_rfc2141 | Urn RFC 2141 String |

### Strings

| Tag | Description |
| - | - |
| alpha | Alpha Only |
| alphanum | Alphanumeric |
| alphanumunicode | Alphanumeric Unicode |
| alphaunicode | Alpha Unicode |
| ascii | ASCII |
| boolean | Boolean |
| contains | Contains |
| containsany | Contains Any |
| containsrune | Contains Rune |
| endsnotwith | Ends With |
| endswith | Ends With |
| excludes | Excludes |
| excludesall | Excludes All |
| excludesrune | Excludes Rune |
| lowercase | Lowercase |
| multibyte | Multi-Byte Characters |
| number | NOT DOCUMENTED IN doc.go |
| numeric | Numeric |
| printascii | Printable ASCII |
| startsnotwith | Starts Not With |
| startswith | Starts With |
| uppercase | Uppercase |

### Format

| Tag | Description |
| - | - |
| base64 | Base64 String |
| base64url | Base64URL String |
| bic | Business Identifier Code (ISO 9362) |
| bcp47_language_tag | Language tag (BCP 47) |
| btc_addr | Bitcoin Address |
| btc_addr_bech32 | Bitcoin Bech32 Address (segwit) |
| datetime | Datetime |
| e164 | e164 formatted phone number |
| email | E-mail String
| eth_addr | Ethereum Address |
| hexadecimal | Hexadecimal String |
| hexcolor | Hexcolor String |
| hsl | HSL String |
| hsla | HSLA String |
| html | HTML Tags |
| html_encoded | HTML Encoded |
| isbn | International Standard Book Number |
| isbn10 | International Standard Book Number 10 |
| isbn13 | International Standard Book Number 13 |
| iso3166_1_alpha2 | Two-letter country code (ISO 3166-1 alpha-2) |
| iso3166_1_alpha3 | Three-letter country code (ISO 3166-1 alpha-3) |
| iso3166_1_alpha_numeric | Numeric country code (ISO 3166-1 numeric) |
| iso3166_2 | Country subdivision code (ISO 3166-2) |
| iso4217 | Currency code (ISO 4217) |
| json | JSON |
| jwt | JSON Web Token (JWT) |
| latitude | Latitude |
| longitude | Longitude |
| postcode_iso3166_alpha2 | Postcode |
| postcode_iso3166_alpha2_field | Postcode |
| rgb | RGB String |
| rgba | RGBA String |
| ssn | Social Security Number SSN |
| timezone | Timezone |
| uuid | Universally Unique Identifier UUID |
| uuid3 | Universally Unique Identifier UUID v3 |
| uuid3_rfc4122 | Universally Unique Identifier UUID v3 RFC4122 |
| uuid4 | Universally Unique Identifier UUID v4 |
| uuid4_rfc4122 | Universally Unique Identifier UUID v4 RFC4122 |
| uuid5 | Universally Unique Identifier UUID v5 |
| uuid5_rfc4122 | Universally Unique Identifier UUID v5 RFC4122 |
| uuid_rfc4122 | Universally Unique Identifier UUID RFC4122 |
| semver | Semantic Versioning 2.0.0 |
| ulid | Universally Unique Lexicographically Sortable Identifier ULID |

### Comparisons

| Tag | Description |
| - | - |
| eq | Equals |
| gt | Greater than|
| gte | Greater than or equal |
| lt | Less Than |
| lte | Less Than or Equal |
| ne | Not Equal |

### Other

| Tag | Description |
| - | - |
| dir | Directory |
| file | File path |
| isdefault | Is Default |
| len | Length |
| max | Maximum |
| min | Minimum |
| oneof | One Of |
| required | Required |
| required_if | Required If |
| required_unless | Required Unless |
| required_with | Required With |
| required_with_all | Required With All |
| required_without | Required Without |
| required_without_all | Required Without All |
| excluded_with | Excluded With |
| excluded_with_all | Excluded With All |
| excluded_without | Excluded Without |
| excluded_without_all | Excluded Without All |
| unique | Unique |

#### Aliases

| Tag | Description |
| - | - |
| iscolor | hexcolor\|rgb\|rgba\|hsl\|hsla |
| country_code | iso3166_1_alpha2\|iso3166_1_alpha3\|iso3166_1_alpha_numeric |
