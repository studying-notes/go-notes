---
date: 2020-07-20T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Redis 客户端操作"  # 文章标题
url:  "posts/go/libraries/tripartite/database/redis"  # 设置网页永久链接
tags: [ "Go", "redis" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [安装](#安装)
- [连接](#连接)
  - [普通连接](#普通连接)
  - [集群模式](#集群模式)
- [基本用法](#基本用法)
  - [Set / Get](#set--get)
  - [ZSet](#zset)
  - [Pipeline](#pipeline)
  - [事务](#事务)
  - [Watch](#watch)
  - [发布订阅](#发布订阅)

## 安装

```shell
go get github.com/go-redis/redis/v8
```

## 连接

### 普通连接

```go
package main

import (
    "context"
    "fmt"

    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ExampleClient() {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    err := rdb.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := rdb.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)

    val2, err := rdb.Get(ctx, "key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exist")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
    }
    // Output: key value
    // key2 does not exist
}
```

### 集群模式

```go
package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{":6371", ":6372", ":6373", ":6374", ":6375", ":6376"},
		Password: "120e204105de1345fda9f27911c02f66",
	})

	ctx, cancel := context.WithCancel(context.Background())

	err := rdb.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})

	if err != nil {
		cancel()
		log.Fatal(err)
	}
}
```

## 基本用法

### Set / Get

```go
func main() {
	rdb.Set("score", 100, 0)
	val, err := rdb.Get("score").Result()
	if err == redis.Nil {
		fmt.Println("not exist")
	} else if err != nil {
		fmt.Println("failed:", err)
	} else {
		fmt.Println("score", val)
	}
}
```

### ZSet

Redis zset 和 set 一样也是 string 类型元素的集合，且不允许重复的成员。不同的是每个元素都会关联一个 double 类型的分数。Redis 正是通过分数来为集合中的成员进行从小到大的排序。zset 的成员是唯一的，但分数却可以重复。

```go
func main() {
	var rdb *redis.Client
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	zsetKey := "language_rank"
	languages := []redis.Z{
		{Score: 90.0, Member: "Golang"},
		{Score: 98.0, Member: "Java"},
		{Score: 95.0, Member: "Python"},
		{Score: 97.0, Member: "JavaScript"},
		{Score: 99.0, Member: "C/C++"},
	}

	// zadd
	num, _ := rdb.ZAdd(zsetKey, languages...).Result()
	fmt.Printf("zadd %d", num)

	// 增加元素权重
	newScore, _ := rdb.ZIncrBy(zsetKey, 10.0, "Golang").Result()
	fmt.Printf("Golang's score is %f now.\n", newScore)

	// 取排名前三
	ret, _ := rdb.ZRevRangeWithScores(zsetKey, 0, 2).Result()
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}

	// 取区间内的
	op := redis.ZRangeBy{Max: "100", Min: "90"}
	ret, _ = rdb.ZRangeByScoreWithScores(zsetKey, op).Result()
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}
}
```

### Pipeline

`Pipeline` 主要是一种网络优化。它本质上意味着客户端缓冲一堆命令并一次性将它们发送到服务器。这些命令不能保证在事务中执行。这样做的好处是节省了每个命令的网络往返时间（RTT）。

```go
func main() {
	pipe := rdb.Pipeline()

	incr := pipe.Incr("pipeline_counter")
	pipe.Expire("pipeline_counter", time.Hour)

	_, err = pipe.Exec()
	fmt.Println(incr.Val(), err)
}
```

上面的代码相当于将以下两个命令一次发给 Redis Server 端执行，与不使用 `Pipeline` 相比能减少一次RTT。下面的代码是等价的：

```go
func main() {
	var incr *redis.IntCmd
	_, err = rdb.Pipelined(func(pipe redis.Pipeliner) error {
		incr = pipe.Incr("pipelined_counter")
		pipe.Expire("pipelined_counter", time.Hour)
		return nil
	})
	fmt.Println(incr.Val(), err)
}
```

### 事务

Redis 是单线程的，因此单个命令始终是原子的，但是来自不同客户端的两个给定命令可以依次执行，例如在它们之间交替执行。但是，`Multi/exec` 能够确保在 `multi/exec` 两个语句之间的命令之间没有其他客户端正在执行命令。

在这种场景我们需要使用 `TxPipeline`。`TxPipeline` 总体上类似于上面的 `Pipeline`，但是它内部会使用 `MULTI/EXEC` 包裹排队的命令。例如：

```go
func main() {
	pipe := rdb.TxPipeline()
	incr := pipe.Incr("tx_pipeline_counter")
	pipe.Expire("tx_pipeline_counter", time.Hour)

	_, err = pipe.Exec()
	fmt.Println(incr.Val(), err)
}
```

类似地：

```go
func main() {
	var incr *redis.IntCmd
	_, err = rdb.TxPipelined(func(pipe redis.Pipeliner) error {
		incr = pipe.Incr("tx_pipelined_counter")
		pipe.Expire("tx_pipelined_counter", time.Hour)
		return nil
	})
	fmt.Println(incr.Val(), err)
}
```

### Watch

在某些场景下，我们除了要使用 `MULTI/EXEC` 命令外，还需要配合使用 `WATCH` 命令。在用户使用 `WATCH` 命令监视某个键之后，直到该用户执行 `EXEC` 命令的这段时间里，如果有其他用户抢先对被监视的键进行了替换、更新、删除等操作，那么当用户尝试执行 `EXEC` 的时候，事务将失败并返回一个错误，用户可以根据这个错误选择重试事务或者放弃事务。

```go
Watch(fn func(*Tx) error, keys ...string) error
```

Watch 方法接收一个函数和一个或多个 key 作为参数。

### 发布订阅

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{":6371", ":6372", ":6373", ":6374", ":6375", ":6376"},
		Password: "120e204105de1345fda9f27911c02f66",
	})

	ctx, cancel := context.WithCancel(context.Background())

	err := rdb.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})

	if err != nil {
		cancel()
		log.Fatal(err)
	}

	subscriber := rdb.Subscribe(ctx, "dev")
	defer subscriber.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-subscriber.Channel():
			log.Println(msg.Payload)
		default:
			err = rdb.Publish(ctx, "dev", "message").Err()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Message published")

			time.Sleep(3 * time.Second)
		}
	}
}
```

```go

```
