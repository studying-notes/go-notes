---
date: 2020-11-16T21:21:50+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 图片处理之 image 标准库"  # 文章标题
url:  "posts/go/io/image/draw"  # 设置网页永久链接
tags: [ "go", "图片处理" ]  # 标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 生成图片

```go
func genImage() {
    m := image.NewRGBA(image.Rect(0, 0, 640, 480))

    draw.Draw(m, m.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

    f, err := os.Create("demo.jpeg")
    if err != nil {
        panic(err)
    }
    err = jpeg.Encode(f, m, nil)
    if err != nil {
        panic(err)
    }
}
```

## 读取图片

```go
func readImage()  {
    f, err := os.Open("ubuntu.png")
    if err != nil {
        panic(err)
    }
    // decode图片
    m, err := png.Decode(f)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%v\n", m.Bounds())     // 图片长宽
    fmt.Printf("%v\n", m.ColorModel()) // 图片颜色模型
    fmt.Printf("%v\n", m.At(100,100))  // 该像素点的颜色
}
```

## 图片裁剪

```go
func readImage() {
    f, err := os.Open("ubuntu.png")
    if err != nil {
        panic(err)
    }
    // decode图片
    m, err := png.Decode(f)
    if err != nil {
        panic(err)
    }
    rgba := m.(*image.RGBA)
    subImage := rgba.SubImage(image.Rect(0, 0, 266, 133)).(*image.RGBA)

    // 保存图片
    create, _ := os.Create("new.png")
    err = png.Encode(create, subImage)
    if err != nil {
        panic(err)
    }
}
```

## 图片转 Base64

```go
func b64() {
    f, err := os.Open("ubuntu.png")
    if err != nil {
        panic(err)
    }
    all, _ := ioutil.ReadAll(f)

    str := base64.StdEncoding.EncodeToString(all)
    
    fmt.Printf("%s\n", str)
}
```

## 图片大小压缩

```go
package main

import (
    "github.com/nfnt/resize"
    "image/jpeg"
    "log"
    "os"
)

func main() {
    // open "test.jpg"
    file, err := os.Open("test.jpg")
    if err != nil {
        log.Fatal(err)
    }
    // decode jpeg into image.Image
    img, err := jpeg.Decode(file)
    if err != nil {
        log.Fatal(err)
    }
    file.Close()

    // resize to width 1000 using Lanczos resampling
    // and preserve aspect ratio
    m := resize.Resize(1000, 0, img, resize.Lanczos3)

    out, err := os.Create("test_resized.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer out.Close()

    // write new image to file
    jpeg.Encode(out, m, nil)
}
```

```go

```

```go

```

