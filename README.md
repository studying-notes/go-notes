# Go 学习笔记

纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。

- 基础知识
    - [Go 语言数组](base/array.md)
    - [Go 语言字符串](base/string.md)
    - [Go 语言切片](base/slice.md)
    - [Go 语言函数](base/func.md)
    - [Go 语言方法](base/method.md)
    - [Go 语言接口](base/interface.md)
    - [Goroutine 内存模型](base/goroutine.md)
    - [常见的并发模式](base/concurrent.md)

- 输入与输出
    - 日志
        - [Log 标准库](io/log/log.md)
        - [Logrus 日志库](io/log/logrus.md)
        - [Zap 日志库](io/log/zap.md)

- [加密与解密](encrypt/README.md)

- 数据库操作
    - [Go 语言操作 MongoDB](database/mongo/mongo.md)
    - [Go 语言操作 Redis](database/redis/redis.md)

    - MySQL
        - [database/sql 标准库](database/mysql/sql.md)
        - [sqlx 拓展库](database/mysql/sqlx.md)

    - 加密的 SQLite 数据库 - SQLCipher
        - [go-sqlcipher 的编译与安装](database/sqlite3/sqlcipher/install.md) - Windows 10 x64
        - [go-sqlcipher 基本用法](database/sqlite3/sqlcipher/usage.md)

- 网络编程
    - Websocket
        - [通过命令行对话](web/websocket/c2s/README.md)
        - [文件修改监视器](web/websocket/watch/README.md)
        - [网页聊天室](web/websocket/chatroom/README.md)
        - [时间回显](web/websocket/echo/README.md)
    - [gRPC](web/grpc/README.md)
