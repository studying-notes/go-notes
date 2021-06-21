package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

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

	var incr *redis.IntCmd
	_, err = rdb.TxPipelined(func(pipe redis.Pipeliner) error {
		incr = pipe.Incr("tx_pipelined_counter")
		pipe.Expire("tx_pipelined_counter", time.Hour)
		return nil
	})
	fmt.Println(incr.Val(), err)
}
