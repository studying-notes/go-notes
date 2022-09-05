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
