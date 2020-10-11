# Kafka

- [Kafka](#kafka)
	- [快捷启动](#快捷启动)
	- [Windows 安装 Kafka](#windows-安装-kafka)
	- [命令行操作](#命令行操作)
		- [创建主题](#创建主题)
		- [查看主题](#查看主题)
		- [创建生产者](#创建生产者)
		- [创建消费者](#创建消费者)
	- [Go 操作 Kafka](#go-操作-kafka)
		- [Conn](#conn)
			- [一个简单的生产者](#一个简单的生产者)
			- [一个简单的消费者](#一个简单的消费者)
		- [Reader](#reader)
			- [Consumer Groups](#consumer-groups)
			- [Explicit Commits](#explicit-commits)
			- [Managing Commits](#managing-commits)
		- [Writer](#writer)
		- [TLS Support](#tls-support)

## 快捷启动

1. 启动 Zookeeper

```powershell
zkserver
```

2. 运行 Kafka

```shell
kafka-server-start "c:/developer/kafka_2.13-2.5.0/config/server.properties"
```

3. 创建主题

```shell
kafka-topics --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
```

4. 查看主题

```shell
kafka-topics --list --zookeeper localhost:2181
```

5. 创建生产者

```shell
kafka-console-producer --broker-list localhost:9092 --topic test
```

6. 创建消费者

```
kafka-console-consumer --bootstrap-server localhost:9092 --topic test --from-beginning
```

## Windows 安装 Kafka

1. 下载 Zookeeper

```
http://archive.apache.org/dist/zookeeper/
```

> 名称里带 `bin` 的是编译好的，不带的是源码！

2. 解压 Zookeeper

```shell
tar -xvzf apache-zookeeper-3.6.1-bin.tar.gz
```

3. 进入解压目录创建编辑 `conf/zoo.cfg` 文件

```ini
# The number of milliseconds of each tick
tickTime=2000
# The number of ticks that the initial 
# synchronization phase can take
initLimit=10
# The number of ticks that can pass between 
# sending a request and getting an acknowledgement
syncLimit=5
# the directory where the snapshot is stored.
# do not use /tmp for storage, /tmp here is just 
# example sakes.
dataDir=c:/developer/data/zookeeper
# the port at which the clients will connect
clientPort=2181
```

4. 启动 Zookeeper

```shell
> zkserver
```

不报错窗口不退出就成功了。

5. 下载 Kafka

```
http://archive.apache.org/dist/kafka/
```

6. 解压 Kafka

```shell
tar -xvzf kafka_2.13-2.5.0.tgz
```

7. 按需求编辑 `config/server.properties` 文件

```ini
zookeeper.connect=localhost:2181
#listeners=PLAINTEXT://:9092
#advertised.listeners=PLAINTEXT://your.host.name:9092
```

8. 运行 Kafka

```shell
> cd C:/Developer/kafka_2.13-2.5.0
> kafka-server-start config/server.properties
```

- `config/server.properties` 必须指定配置文件，运行后窗口不退出。

## 命令行操作

### 创建主题

```shell
> kafka-topics --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
Created topic test.
```

### 查看主题

```shell
> kafka-topics --list --zookeeper localhost:2181
test
```

### 创建生产者

```shell
> kafka-console-producer --broker-list localhost:9092 --topic test
>
```

### 创建消费者

```
> kafka-console-consumer --bootstrap-server localhost:9092 --topic test --from-beginning
```

就是一个接受窗口，在生产者窗口发送消息就会出现在消费者窗口。

## Go 操作 Kafka

```shell
go get github.com/segmentio/kafka-go
```

### Conn

这个是比较低级的接口。

#### 一个简单的生产者

```go
func main() {
	topic := "test"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp",
		"localhost:9092", topic, partition)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, _ = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
}
```

#### 一个简单的消费者

```go
func main() {
	topic := "test"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp",
		"localhost:9092", topic, partition)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6)
	defer batch.Close()
	b := make([]byte, 10e3)
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b))
	}
}
```

### Reader

Reader 旨在简化从单个主题分区对进行消费的典型用例。Reader 还会自动处理重新连接和偏移量管理，并公开使用 Go contexts 支持异步取消和超时的 API。

```go
func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "test",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
	defer r.Close()

	_ = r.SetOffset(42)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}
```

#### Consumer Groups

支持 Kafka 消费者组，指定 `GroupID` 即可，`ReadMessage` 自动提交 `offsets`。

```go
func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		GroupID:  "test-id",
		Topic:    "test",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()

	_ = r.SetOffset(42)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
```

- `(*Reader).SetOffset` will return an error when GroupID is set
- `(*Reader).Offset` will always return `-1` when GroupID is set
- `(*Reader).Lag` will always return `-1` when GroupID is set
- `(*Reader).ReadLag` will return an error when GroupID is set
- `(*Reader).Stats` will return a partition of `-1` when GroupID is set

#### Explicit Commits

```go
ctx := context.Background()
for {
m, err := r.FetchMessage(ctx)
if err != nil {
	break
}
fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
_ = r.CommitMessages(ctx, m)
}
```

#### Managing Commits

默认 `CommitMessages` 会将偏移量同步提交给 Kafka。 为了提高性能，可以通过在 `ReaderConfig` 上设置 `CommitInterval` 来定期向 Kafka 提交偏移量。

```go
r := kafka.NewReader(kafka.ReaderConfig{
    Brokers:        []string{"localhost:9092"},
    GroupID:        "consumer-group-id",
    Topic:          "topic-A",
    MinBytes:       10e3, // 10KB
    MaxBytes:       10e6, // 10MB
    CommitInterval: time.Second, // flushes commits to Kafka every second
})
```

### Writer

Writer 是比 Conn 更高级的接口，提供了附加特性：

- Automatic retries and reconnections on errors.
- Configurable distribution of messages across available partitions.
- Synchronous or asynchronous writes of messages to Kafka.
- Asynchronous cancellation using contexts.
- Flushing of pending messages on close to support graceful shutdowns.

```go
func main() {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "test",
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	_ = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
}
```

> `kafka.Message` 拥有 `Topic` 和 `Partition` 字段，但是在写消息的时候不可以设置，这两个字段仅仅是为了方便阅读。

### TLS Support

```go
dialer := &kafka.Dialer{
    Timeout:   10 * time.Second,
    DualStack: true,
    TLS:       &tls.Config{...tls config...},
}

conn, err := dialer.DialContext(ctx, "tcp", "localhost:9093")
```

```go
dialer := &kafka.Dialer{
    Timeout:   10 * time.Second,
    DualStack: true,
    TLS:       &tls.Config{...tls config...},
}

r := kafka.NewReader(kafka.ReaderConfig{
    Brokers:        []string{"localhost:9093"},
    GroupID:        "consumer-group-id",
    Topic:          "topic-A",
    Dialer:         dialer,
})
```

```go
dialer := &kafka.Dialer{
    Timeout:   10 * time.Second,
    DualStack: true,
    TLS:       &tls.Config{...tls config...},
}

w := kafka.NewWriter(kafka.WriterConfig{
	Brokers: []string{"localhost:9093"},
	Topic:   "topic-A",
	Balancer: &kafka.Hash{},
	Dialer:   dialer,
})
```
