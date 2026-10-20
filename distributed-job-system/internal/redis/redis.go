package redis

import (
	"context"
	goredis "github.com/redis/go-redis/v9"
)

var Client *goredis.Client

func NewClient(addr string) {
	Client = goredis.NewClient(&goredis.Options{
		Addr: addr,
	})
}

func Ping() error {
	cxt := context.Background()

	return Client.Ping(cxt).Err()
}
