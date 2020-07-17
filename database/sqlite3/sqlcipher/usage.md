# go-sqlcipher 基本用法

## 创建被保护的数据库

```go
package main

import (
	"database/sql"
	_ "github.com/xeodou/go-sqlcipher"
)

func func main() {
    sql.Open("sqlite3", databasefile +"?_key=password")
}
```

经测试，现在只支持自定义密钥 `key`，其他设置了都无效，最终都会变成默认设置。

## 打开被保护的数据库文件

折腾半天，走了很多错误的方向，一是 Windows 平台找不到编译好的最新版本，找了个旧版本（当时还没发现这个问题），怎么也打不开生成的文件；二是根据 StackOverflow 的回答找了一个专门打开 SQLite 的可视化工具（SQLiteStudio），然而必须自己输除密钥外的其他设置，怎么输都错。于是我搜索寻找默认设置，然而从官方文档到源代码，只找到一丝眉目。

几近绝望，最终在一篇讲微信加密数据库的博客里发现了一款工具：

```
https://download.sqlitebrowser.org/DB.Browser.for.SQLite-3.12.0-win64.msi
```

提供了新旧版本的默认设置，问题解决了。
