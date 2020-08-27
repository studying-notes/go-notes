# Go 反射

- [Go 反射](#go-反射)
	- [遍历结构体](#遍历结构体)
		- [获取标签值](#获取标签值)
		- [获取字段值](#获取字段值)
		- [处理字段值](#处理字段值)

## 遍历结构体

### 获取标签值

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

### 获取字段值

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

### 处理字段值

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
				for i := 0; i < obj.Len(); i++ {
					_ = obj.Index(i).String() // 提取切片数据
				}
			}
		}
	}
}
```
