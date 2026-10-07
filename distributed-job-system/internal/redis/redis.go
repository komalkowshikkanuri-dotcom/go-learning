package redis

import (
	"context"
	goredis "github.com/redis/go-redis/v9"
)

var Client *goredis.Client

func NewClient() {
	Client = goredis.NewClient(&goredis.Options{
		Addr: "localhost:6379",
	})
}

func Ping() error {
	cxt := context.Background()

	return Client.Ping(cxt).Err()
}
