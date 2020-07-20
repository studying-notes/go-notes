# Go 语言操作 MongoDB

```
https://www.mongodb.com/download-center/community
```

## 基本用法

### 服务端启动

```shell
mongod --dbpath="c:/data
```

### 客户端启动

```shell
mongo
```

### 数据库常用命令

**查看数据库**

```shell
> show dbs
admin   0.000GB
config  0.000GB
local   0.000GB
```

**切换到指定数据库**

```shell
> use dev
switched to db dev
```

**显示当前所在数据库**

```shell
> db
dev 
```

**删除当前数据库**

```shell
> db.dropDatabase()
{ "ok" : 1 }
```

### 数据集常用命令

**创建数据集** - `db.createCollection(name,options)`

- `name`：数据集名称；
- `options`：可选参数，指定内存大小和索引。

```shell
> db.createCollection("student")
{ "ok" : 1 }
```

**查看当前数据库中所有集合**

```shell
> show collections
student
```

**删除指定数据集**

```shell
> db.student.drop()
true
```

### 文档常用命令

#### 插入一条文档

```shell
> db.student.insertOne({name:"小王子",age:18})
{
        "acknowledged" : true,
        "insertedId" : ObjectId("5f15626f67934a3e3a9dc597")
}
```

#### 插入多条文档

```shell
> db.student.insertMany([{name:"张三", age:20}, {name:"李四", age:25}])
{
        "acknowledged" : true,
        "insertedIds" : [
                ObjectId("5f1562e067934a3e3a9dc598"),
                ObjectId("5f1562e067934a3e3a9dc599")
        ]
}
```

#### 查询所有文档

```shell
> db.student.find()
{ "_id" : ObjectId("5f15626f67934a3e3a9dc597"), "name" : "小王子", "age" : 18 }
{ "_id" : ObjectId("5f1562e067934a3e3a9dc598"), "name" : "张三", "age" : 20 }
{ "_id" : ObjectId("5f1562e067934a3e3a9dc599"), "name" : "李四", "age" : 25 }
```

#### 条件查询

```shell
> db.student.find({age:{$gt:20}})
{ "_id" : ObjectId("5f1562e067934a3e3a9dc599"), "name" : "李四", "age" : 25 }
```

#### 更新文档

```shell
> db.student.update({name: "小王子"}, {name: "老王子", age: 98})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })

> db.student.find()
{ "_id" : ObjectId("5f15626f67934a3e3a9dc597"), "name" : "老王子", "age" : 98 }
{ "_id" : ObjectId("5f1562e067934a3e3a9dc598"), "name" : "张三", "age" : 20 }
{ "_id" : ObjectId("5f1562e067934a3e3a9dc599"), "name" : "李四", "age" : 25 }
```

#### 删除文档

```shell
> db.student.deleteOne({name:"李四"}) 
{ "acknowledged" : true, "deletedCount" : 1 }
```

## Go 语言操作 MongoDB

```shell
go get go.mongodb.org/mongo-driver/mongo
```

### 通过 Go 代码连接 MongoDB

```go
package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB.")
}
```

处理数据库中的数据集：

```go
collection := client.Database("dev").Collection("student")
```

断开与 MongoDB 的连接：

```go
_ = client.Disconnect(context.TODO())
```

### 连接池模式

```go
func main() {
	var timeout time.Duration = time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions.SetMaxPoolSize(64)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	client.Database("dev")
}
```

### BSON

MongoDB 中的 JSON 文档存储在名为 BSON 的二进制表示中。与其他将 JSON 数据存储为简单字符串和数字的数据库不同，BSON 编码扩展了 JSON 表示，使其包含额外的类型，如 int、long、date、浮点数和 decimal128。这使得应用程序更容易可靠地处理、排序和比较数据。

连接 MongoDB 的 Go 驱动程序中有两大类型表示 BSON 数据：`D` 和 `Raw`。

类型 `D` 家族被用来简洁地构建使用本地 Go 类型的 BSON 对象。这对于构造传递给 MongoDB 的命令特别有用。`D` 家族包括四类:

- D：一个 BSON 文档。这种类型应该在顺序重要的情况下使用，比如 MongoDB 命令；
- M：一张无序的 Map。它和 D 是一样的，只是它不保持顺序；
- A：一个 BSON 数组；
- E：D 里面的一个元素。

```go
import "go.mongodb.org/mongo-driver/bson"
```

### CRUD

```go
type Student struct {
	Name string
	Age int
}

s1 := Student{"小红", 12}
s2 := Student{"小兰", 10}
s3 := Student{"小黄", 11}
```

#### 插入文档

**插入一条文档记录**

```go
func main() {
	collection := client.Database("dev").Collection("student")
	insertResult, err := collection.InsertOne(context.TODO(), s1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}
```

**插入多条文档记录**

```go
func main() {
    collection := client.Database("dev").Collection("student")
	students := []interface{}{s2, s3}
	insertManyResult, err := collection.InsertMany(context.TODO(), students)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}
```

#### 更新文档

```go
func main() {
    collection := client.Database("dev").Collection("student")
	filter := bson.D{{"name", "小兰"}}
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)
}
```

#### 查找文档

查找一个文档，需要一个 `filter` 文档，以及一个指向可以将结果解码为其值的指针，要查找单个文档，使用 `collection.FindOne()`，这个方法返回一个可以解码为值的结果。

```go
func main() {
	collection := client.Database("dev").Collection("student")
	filter := bson.D{{"name", "小兰"}}
	var result Student
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)
}
```

查找多个文档，可用 `collection.Find()`。此方法返回一个游标，游标提供了一个文档流，可以通过它一次迭代解码一个文档。当游标用完之后，应该关闭游标。下面的示例将使用 `options` 包设置一个限制以便只返回两个文档。

```go
func main() {
	collection := client.Database("dev").Collection("student")
	
	findOptions := options.Find()
	findOptions.SetLimit(2)
	var results []*Student
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Student
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	_ = cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
}
```

#### 删除文档

可以使用 `collection.DeleteOne()` 或 `collection.DeleteMany()` 删除文档。如果传递 `bson.D{{}}` 作为过滤器参数，它将匹配数据集中的所有文档。还可以使用 `collection. drop()` 删除整个数据集。

```go
func main() {
	collection := client.Database("dev").Collection("student")

	// 删除名字是小黄的那个
	deleteResult1, err := collection.DeleteOne(context.TODO(), bson.D{{"name", "小黄"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult1.DeletedCount)

	// 删除所有
	deleteResult2, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult2.DeletedCount)
}
```
