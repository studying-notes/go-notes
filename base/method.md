# Go 语言方法

方法一般是面向对象编程（Object-Oriented Programming，OOP）的一个特性，Go 语言的方法是关联到类型的，这样可以在编译阶段完成方法的静态绑定。一个面向对象的程序会用方法来表达其属性对应的操作，这样使用这个对象的用户就不需要直接去操作对象，而是借助方法来做这些事情。

面向对象编程更多的只是一种思想，很多号称支持面向对象编程的语言只是将经常用到的特性内置到语言中了而已。Go 语言的祖先 C 语言虽然不是一个支持面向对象的语言，但是 C 语言的标准库中的 File 相关的函数也用到了面向对象编程的思想。下面我们实现一组 C 语言风格的 File 函数：

```go
// 文件对象
type File struct {
    fd int
}

// 打开文件
func OpenFile(name string) (f *File, err error) {
    // ...
}

// 关闭文件
func CloseFile(f *File) error {
    // ...
}

// 读文件数据
func ReadFile(f *File, int64 offset, data []byte) int {
    // ...
}
```

其中 `OpenFile()` 类似于构造函数，用于打开文件对象，`CloseFile()` 类似于析构函数，用于关闭文件对象，`ReadFile()` 则类似于普通的成员函数，这 3 个函数都是普通函数。`CloseFile()` 和 `ReadFile()` 作为普通函数，需要占用包级空间中的名字资源。不过 `CloseFile()` 和 `ReadFile()` 函数只是针对 File 类型对象的操作，这时候我们更希望这类函数和操作对象的类型紧密绑定在一起。

Go 语言中的做法是将函数 `CloseFile()` 和 `ReadFile()` 的第一个参数移动到函数名的开头：

```go
// 关闭文件
func (f *File) CloseFile() error {
    // ...
}

// 读文件数据
func (f *File) ReadFile(int64 offset, data []byte) int {
    // ...
}
```

这样的话，函数 `CloseFile()` 和 `ReadFile()` 就成了 `File` 类型独有的方法了（而不是 `File` 对象方法）。它们也不再占用包级空间中的名字资源，同时 `File` 类型已经明确了它们的操作对象，因此方法名字一般简化为 `Close` 和 `Read`：

```go
// 关闭文件
func (f *File) Close() error {
    // ...
}

// 读文件数据
func (f *File) Read(int64 offset, data []byte) int {
    // ...
}
```

我们可以给任何自定义类型添加一个或多个方法。**每种类型对应的方法必须和类型的定义在同一个包中**，因此是无法给内置类型添加方法的（因为方法的定义和类型的定义不在一个包中）。对于给定的类型，每个方法的名字必须是唯一的，同时**方法和函数一样也不支持重载**。

方法由函数演变而来，只是将函数的第一个对象参数移动到了函数名前面了而已。因此我们依然可以按照原始的过程式思维来使用方法。通过称为**方法表达式**的特性可以**将方法还原为普通类型的函数**：

```go
// 不依赖具体的文件对象
// func CloseFile(f *File) error
var CloseFile = (*File).Close

// 不依赖具体的文件对象
// func ReadFile(f *File, int64 offset, data []byte) int
var ReadFile = (*File).Read

// 文件处理
f, _ := OpenFile("foo.dat")
ReadFile(f, 0, data)
CloseFile(f)
```

在有些场景更关心一组相似的操作。例如，`Read()` 读取一些数组，然后调用 `Close()` 关闭。此时的环境中，用户并不关心操作对象的类型，只要能满足通用的 `Read()` 和 `Close()` 行为就可以了。不过在方法表达式中，因为得到的 `ReadFile()` 和 `CloseFile()` 函数参数中含有 `File` 这个特有的类型参数，这使得 `File` 相关的方法无法与其他不是 `File` 类型但是有着相同 `Read()` 和 `Close()` 方法的对象无缝适配。我们可以通过结合闭包特性来消除方法表达式中第一个参数类型的差异：

```go
// 先打开文件对象
f, _ := OpenFile("foo.dat")

// 绑定到了f对象
// func Close() error
var Close = func Close() error {
    return (*File).Close(f)
}

// 绑定到了f对象
// func Read(int64 offset, data []byte) int
var Read = func Read(int64 offset, data []byte) int {
    return (*File).Read(f, offset, data)
}

// 文件处理
Read(0, data)
Close()
```

这刚好是方法值也要解决的问题。我们用方法值特性可以简化实现：

```go
// 先打开文件对象
f, _ := OpenFile("foo.dat")

// 方法值：绑定到了f对象
// func Close() error
var Close = f.Close

// 方法值：绑定到了f对象
// func Read(int64 offset, data []byte) int
var Read = f.Read

// 文件处理
Read(0, data)
Close()
```

Go 语言不支持传统面向对象中的继承特性，而是以自己特有的组合方式支持了方法的继承。Go 语言中，通过在结构体内置匿名的成员来实现继承：

```go
import "image/color"

type Point struct{ X, Y float64 }

type ColoredPoint struct {
	Point
	Color color.RGBA
}
```

虽然我们可以将 `ColoredPoint` 定义为一个有 3 个字段的扁平结构的结构体，但是这里将 `Point` 嵌入 `ColoredPoint` 来提供 X 和 Y 这两个字段：

```go
var cp ColoredPoint
cp.X = 1
fmt.Println(cp.Point.X) // "1"
cp.Point.Y = 2
fmt.Println(cp.Y) // "2"
```

通过嵌入匿名的成员，不仅可以继承匿名成员的内部成员，而且可以继承匿名成员类型所对应的方法。我们一般会将 `Point` 看作基类，把 `ColoredPoint` 看作 `Point` 的继承类或子类。不过这种方式继承的方法并不能实现多态特性。所有继承来的方法的接收者参数依然是那个匿名成员本身，而不是当前的变量。

```go
type Cache struct {
    m map[string]string
    sync.Mutex
}

func (p *Cache) Lookup(key string) string {
    p.Lock()
    defer p.Unlock()

    return p.m[key]
}
```

`Cache` 结构体类型通过嵌入一个匿名的 `sync.Mutex` 来继承它的方法 `Lock()` 和 `Unlock()`。但是在调用 `p.Lock()` 和 `p.Unlock()` 时，`p` 并不是方法 `Lock()` 和 `Unlock()` 的真正接收者，而是会将它们展开为 `p.Mutex.Lock()` 和 `p.Mutex.Unlock()` 调用。这种展开是编译期完成的，并没有运行时代价。

在传统的面向对象语言的继承中，子类的方法是在运行时动态绑定到对象的，因此基类实现的某些方法看到的 `this` 可能不是基类类型对应的对象，这个特性会导致基类方法运行的不确定性。而在 Go 语言通过嵌入匿名的成员来“继承”的基类方法，`this` 就是实现该方法的类型的对象，Go 语言中方法是编译时静态绑定的。如果需要虚函数的多态特性，我们需要借助  语言接口来实现。
