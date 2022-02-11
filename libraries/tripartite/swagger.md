---
date: 2020-08-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "swagger - 通过注释在框架中集成 Swagger"  # 文章标题
url:  "posts/go/libraries/tripartite/swagger"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go", "swagger" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

> 不支持 map 和 array，别折腾，无必要

- [Swagger Editor](#swagger-editor)
- [安装 swaggo](#安装-swaggo)
- [注释基本信息](#注释基本信息)
- [注释接口信息](#注释接口信息)
- [常用请求示例](#常用请求示例)
  - [Headers 鉴权](#headers-鉴权)
  - [multipart/form-data](#multipartform-data)
  - [application/x-www-form-urlencoded](#applicationx-www-form-urlencoded)
  - [多个路径参数](#多个路径参数)
  - [表单上传文件](#表单上传文件)
  - [指明参数属性格式](#指明参数属性格式)
- [常用响应示例](#常用响应示例)
  - [数组类型数据](#数组类型数据)
  - [在注释中组合结构体](#在注释中组合结构体)
  - [添加 Headers](#添加-headers)
  - [给出示例值](#给出示例值)
  - [给出字段描述](#给出字段描述)

## Swagger Editor

建议用 Docker 启动：

```shell
docker pull swaggerapi/swagger-editor
docker run --restart=always -d -p 8080:8080 swaggerapi/swagger-editor
```

```shell
localhost:8080
```

## 安装 swaggo

可能还是找不到命令：

https://github.com/swaggo/swag/issues/209

```shell
go get -u github.com/swaggo/swag/cmd/swag

go install github.com/swaggo/swag/cmd/swag@latest
```

在安装完 Swagger 关联库后，就需要项目里的 API 接口编写注解，以便后续在生成时能够正确地运行。

![](https://i.loli.net/2021/04/03/KOUCvbT46YXdwHf.png)

最后执行以下命令生成相关文档：

```shell
swag init
```

## 注释基本信息

基本信息必须在 main.go 的 main 函数上添加注释，或者 `-g` 参数指定文件。

```shell
swag init -g http/api.go
```

```go
// @title Gin Swagger
// @version 1.0
// @description Gin Swagger 示例项目
// @termsOfService https://github.com/fujiawei-dev

// @contact.name Rustle Karl
// @contact.url https://github.com/fujiawei-dev
// @contact.email fu.jiawei@outlook.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @basePath /api/v1/

// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization


func main() {
 // 省略其他代码
}
```

## 注释接口信息

| 注释属性 | 描述 |
| ----------- | ----------- |
| description | 描述 |
| id | 全局唯一标识符 |
| tags | 接口标注 |
| summary | 简述 |
| accept | 请求类型 |
| produce | 响应类型 |
| param | 参数 `param name`,`param type`,`data type`,`is required`,`comment` `attribute(optional)` |
| success | 请求成功后返回 `return code`,`{param type}`,`data type`,`comment` |
| failure | 请求失败后返回 `return code`,`{param type}`,`data type`,`comment` |
| router | 请求路由及请求方式 `path`,`[httpMethod]` |

| 数据类型 | 注释可填 |
| ----------- | ----------- |
| application/json | application/json, json |
| text/xml | text/xml, xml |
| text/plain | text/plain, plain |
| html | text/html, html |
| multipart/form-data | multipart/form-data, mpfd |
| application/x-www-form-urlencoded | application/x-www-form-urlencoded, x-www-form-urlencoded |
| application/vnd.api+json | application/vnd.api+json, json-api |
| application/x-json-stream | application/x-json-stream, json-stream |
| application/octet-stream | application/octet-stream, octet-stream |
| image/png | image/png, png |
| image/jpeg | image/jpeg, jpeg |
| image/gif | image/gif, gif |

## 常用请求示例

### Headers 鉴权

```go
// @Param Authorization header string true "鉴权"
```

### multipart/form-data

```go
// @Accept mpfd
// @Param app_key formData string true "账号"
// @Param app_secret formData string true "密码"
```

```shell
curl -X POST "http://localhost:8000/auth" -H "accept: application/json" -H "Content-Type: multipart/form-data" -F "app_key=admin" -F "app_secret=admin"
```

### application/x-www-form-urlencoded

```go
// @Accept application/x-www-form-urlencoded

// @Param app_key formData string true "账号"
// @Param app_secret formData string true "密码"

// @Param username formData string true "用户名"
// @Param password formData string true "密码"
```

```shell
curl -X POST "http://localhost:8000/auth" -H "accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -d "app_key=admin&app_secret=admin"
```

一般情况下，Gin 两者都可以解析。

### 多个路径参数

```go
// @Param group_id path int true "Group ID"
// @Param account_id path int true "Account ID"
// @Router /examples/groups/{group_id}/accounts/{account_id} [get]
```

### 表单上传文件

```go
// @Accept mpfd
// @Param file formData file true "上传文件"
```

### 指明参数属性格式

```go
// @Param q query string false "email" Format(email)
// @Param link formData string true "url" Format(uri)
// @Param enumstring query string false "string enums" Enums(A, B, C)
// @Param enumint query int false "int enums" Enums(1, 2, 3)
// @Param enumnumber query number false "int enums" Enums(1.1, 1.2, 1.3)
// @Param string query string false "string valid" minlength(5) maxlength(10)
// @Param int query int false "int valid" minimum(1) maximum(10)
// @Param default query string false "string default" default(A)
// @Param collection query []string false "string collection" collectionFormat(multi)
```

## 常用响应示例

### 数组类型数据

```go
// @Success 200 {array} response.Bottle
```

### 在注释中组合结构体

```go
type JSONResult struct {
    Code    int          `json:"code" `
    Message string       `json:"message"`
    Data    interface{}  `json:"data"`
}

type Order struct { //in `proto` package
    Id  uint            `json:"id"`
    Data  interface{}   `json:"data"`
}

// JSONResult's data field will be overridden by the specific type proto.Order
// @success 200 {object} jsonresult.JSONResult{data=proto.Order} "desc"
// @success 200 {object} jsonresult.JSONResult{data=string} "desc"
// @success 200 {object} jsonresult.JSONResult{data=[]string} "desc"
```

### 添加 Headers

```go
// @Success 200 {string} string	"ok"
// @Header 200 {string} Location "/entity"
// @Header 200 {string} Token "qwerty"
```

### 给出示例值

```go
type Account struct {
    ID   int    `json:"id" example:"1"`
    Name string `json:"name" example:"account"`
    PhotoUrls []string `json:"photo_urls" example:"http://image/1.jpg,http://image/2.jpg"`
}
```

### 给出字段描述

```go
type Account struct {
	// ID this is userid
	ID   int    `json:"id"`
	Name string `json:"name"` // This is Name
}
```
