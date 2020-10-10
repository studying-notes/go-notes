# Go 的反射机制

- [Go 的反射机制](#go-的反射机制)
	- [反射函数和类型](#反射函数和类型)
		- [Type](#type)
		- [Value](#value)
	- [遍历结构体](#遍历结构体)
		- [获取标签值](#获取标签值)
		- [获取字段值](#获取字段值)
		- [处理字段值](#处理字段值)

## 反射函数和类型

reflect 实现了运行时的反射能力，能够让程序操作不同类型的对象。

- `reflect.TypeOf` 获取类型信息，返回 `Type` 类型；
- `reflect.ValueOf` 获取数据的运行时表示，返回 `Value` 类型。

### Type

类型 `Type` 是反射包定义的一个接口，我们可以使用 reflect.TypeOf 函数获取任意变量的的类型，`Type` 接口中定义了多种方法，比如 `MethodByName` 可以获取当前类型对应方法的引用、`Implements` 可以判断当前类型是否实现了某个接口：

```go
type Type interface {
        Align() int
        FieldAlign() int
        Method(int) Method
        MethodByName(string) (Method, bool)
        NumMethod() int
        // ...
        Implements(u Type) bool
        // ...
}
```

### Value

`Value` 的类型与 `Type` 不同，它被声明成了结构体。这个结构体没有对外暴露的字段，但是提供了获取或者写入数据的方法：


```go
type Value struct {
        // ...
}

func (v Value) Addr() Value
func (v Value) Bool() bool
func (v Value) Bytes() []byte
// ...
```

反射中的所有方法基本都是围绕着 `Type` 和 `Value` 这两个类型设计的。通过 `reflect.TypeOf`、`reflect.ValueOf` 可以将一个普通的变量转换成 `Type` 和 `Value`，随后就可以使用相关方法对它们进行复杂的操作。

```go

```

```go

```

```go

```

```go

```

```go

```



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
