---
date: 2022-09-04T16:23:46+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 通过反射获取标签"  # 文章标题
url:  "posts/go/libraries/standard/reflect/tag"  # 设置网页永久链接
tags: [ "Go", "tag" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 获取标签值

提取非嵌套结构体指定标签的值

```go
func ExtractTagValue(i interface{}, tag string) (tagValues []string) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && i != nil {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		types := v.Type()
		for i := 0; i < v.NumField(); i++ {
			tagValues = append(tagValues, types.Field(i).Tag.Get(tag))
		}
	}
	return tagValues
}
```

```go
type Fruit struct {
	ID    string   `json:"id"`
	Name  []string `json:"name"`
	Price string   `json:"price"`
	Area  `json:"area"`
}

func main() {
	fmt.Println(ExtractTagValue(Fruit{}, "json"))
}
```

```
[id name price area]
```

## 获取字段值

获取结构体字段的值

```go
func ExtractFieldValue(i interface{}) (fieldValues []interface{}) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && i != nil {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			fieldValues = append(fieldValues, f.Interface())
		}
	}
	return fieldValues
}
```

```go
func main() {
	fruit := Fruit{ID: "1", Name: []string{"apple", "nut"}, Price: "12",
		Area: Area{Length: "20", Width: "30"}}
	fmt.Println(ExtractFieldValue(fruit))
}
```

```
[1 [apple nut] 12 {20 30}]
```

### 修改字段值

修改结构体字段值

```go
func ModifyFieldValue(ptr interface{}, handle func(string) string) {
	types := reflect.TypeOf(ptr)
	values := reflect.ValueOf(ptr)
	if types.Kind() != reflect.Ptr {
		return // 必须传入指针才能修改原结构体
	}
	types = types.Elem()
	values = values.Elem()
	if values.Kind() == reflect.Struct {
		for i := 0; i < types.NumField(); i++ {
			f := values.Field(i)
			switch f.Kind() {
			case reflect.String: // 字符串类型
				f.Set(reflect.ValueOf(handle(f.String()))) // 设置新字段值
			case reflect.Slice:
				obj := reflect.ValueOf(f.Interface())
				for j := 0; j < obj.Len(); j++ {
					_ = obj.Index(j).String() // 提取切片数据
				}
			}
		}
	}
}
```

```go
func main() {
	ModifyFieldValue(&fruit, func(s string) string {
		return "modify-" + s
	})
	fmt.Printf("%+v\n", fruit)
}
```

```
{ID:modify-1 Name:[apple nut] Price:modify-12 Area:{Length:20 Width:30}}
```

```go

```
