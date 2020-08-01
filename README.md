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
    - [错误和异常](base/error.md)
    - [断言与类型转换](base/assert.md)
    - [Go 语言常见的坑](base/note.md)

- [算法与数据结构](algorithm/README.md)
  - [链表](algorithm/struct/link/README.md)
  - [栈](algorithm/struct/stack/README.md)

- 输入与输出
    - 日志
        - [Log 标准库](io/log/log.md)
        - [Logrus 日志库](io/log/logrus.md)
        - [Zap 日志库](io/log/zap.md)

- 文本操作
    - 字符串
        - [strconv](strings/strconv.md) - 字符串与其他类型相互转换

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
    - [Message Queue](web/mq/README.md)
        - [NSQ](web/mq/nsq/README.md) - Go 语言编写的一个开源的实时分布式内存消息队列
        - [RabbitMQ](web/mq/rabbitmq/README.md) - Erlang 语言编写，一个重量级，适合于企业级的开发
        - [Kafka](web/mq/kafka/README.md) - Jave 语言编写，一个高性能跨语言分布式发布/订阅消息队列系统
    - Websocket
        - [通过命令行对话](web/websocket/c2s/README.md)
        - [文件修改监视器](web/websocket/watch/README.md)
        - [网页聊天室](web/websocket/chatroom/README.md)
        - [时间回显](web/websocket/echo/README.md)
    - [gRPC](web/grpc/README.md)
