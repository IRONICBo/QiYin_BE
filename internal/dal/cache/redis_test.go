package cache

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestSub(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "qiyin",
		DB:       0,
	})

	ctx := context.Background()
	pubsub := client.PSubscribe(ctx, "__keyevent@0__:expired")

	ch := pubsub.Channel()
	go func() {
		for msg := range ch {
			fmt.Println("Events:", msg.Payload)
		}
	}()

	err := client.Set(ctx, "mykey", "myvalue", 5*time.Second).Err()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	err = pubsub.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
