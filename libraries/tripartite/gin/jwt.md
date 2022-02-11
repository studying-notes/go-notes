---
date: 2020-07-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "JWT 鉴权"  # 文章标题
url:  "posts/gin/project/reload"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

## 简介

```
go get -u github.com/dgrijalva/jwt-go
```

JSON Web 令牌（JWT）是一个开放标准（RFC7519），它定义了一种紧凑且自包含的方式，用于在各方之间以 JSON 对象安全地传输信息。由于此信息是经过数字签名的，因此可以被验证和信任。我们可以使用 RSA 或 ECDSA 的公用或专用密钥对 JWT 进行签名。

![](../imgs/jwt.png)

## 组成部分

JWT 是以紧凑的形式由三部分组成的，这三部分之间以点“.”分隔，组成“xxxxx.yyyyy.zzzzz”的格式，三个部分的含义如下：

- Header：头部。
- Payload：有效载荷。
- Signature：签名。

### Header

Header （头部）通常由两部分组成，分别是令牌的类型和所使用的签名算法（HMAC SHA256、RSA等），它们会组成一个JSON对象，用于描述其元数据。例如：

```json
{
    "alg": "HS256",
    "typ": "JWT"
}
```

在上述 JSON 对象中，alg 字段用来表示使用的签名算法，默认是 HMAC SHA256（HS256）。typ 字段用来表示使用的令牌类型，这里使用的是 JWT。最后，用 base64UrlEncode 算法对上面的 JSON 对象进行转换，使其成为 JWT 的第一部分。

### Payload

Payload（有效负载）是一个JSON对象，主要用于存储在JWT中实际传输的数据。例如：

```json
{
    "sub": "1234567890",
    "name": "rustle",
    "admin": "true"
}
```

- aud（Audience）：受众，即接受 JWT 的一方。
- exp（ExpiresAt）：所签发的 JWT 过期时间，过期时间必须大于签发时间。
- jti（JWT Id）：JWT的 唯一标识。
- iat（IssuedAt）：签发时间
- iss（Issuer）：JWT 的签发者。
- nbf（Not Before）：JWT 的生效时间，如果未到这个时间，则不可用。
- sub（Subject）：主题。

同样，使用 base64UrlEncode 算法对该 JSON 对象进行转换，使其成为 JWT Token 的第二部分。
需要注意的是，JWT 在转换时用的是 base64UrlEncode 算法，而该算法是可逆的，因此一些敏感信息建议不要放到 JWT 中。如果一定要放，则应进行一定的加密处理。

### Signature

Signature（签名）部分是对前面两个部分（Header+Payload）进行约定算法和规则的签名。签名一般用于校验消息在整个过程中有没有被篡改，并且对使用了私钥进行签名的令牌，它还可以验证 JWT 的发送者是否是它的真实身份。

### Base64UrlEncode 算法

Base64UrlEncode 算法是 Base64 算法的变种。为什么要变呢？原因是 JWT 令牌经常被放在 Header 或 Query Param 中，即 URL 中。
而在 URL 中，一些个别字符是有特殊意义的，如“+”“/”“=”等。因此在 Base64 UrlEncode 算法中，会对其进行替换。例如，把“+”替换为“-”、把“/”替换为“_”，而″=″ 会被忽略处理，以此保证 JWT 令牌在 URL 中的可用性和准确性。

## 使用场景

首先，在内部约定好 JWT 令牌的交流方式，比如可以存储在 Header、Query Param、cookie 或 session 中，最常见的是存储在 Header 中。然后，服务器端提供一个获取 JWT 令牌的接口方法，返回给客户端使用。当客户端请求其余接口时，需要带上所签发的 JWT 令牌，而服务器端接口也会到约定位置获取 JWT 令牌进行鉴权处理。
