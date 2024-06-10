package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Client struct {
	Client *redis.Client
}

func NewRedisClient(viper *viper.Viper) *Client {
	host := viper.GetString("REDIS_HOST")
	port := viper.GetString("REDIS_PORT")
	//password := viper.GetString("REDIS_PASSWORD")
	addr := fmt.Sprintf("%s:%s", host, port)

	options := &redis.Options{
		Addr: addr,
		//Password: password,
		DB: 0, // use default DB
	}

	client := redis.NewClient(options)
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	log.Println("Connected to Redis")
	return &Client{Client: client}
}

func (r *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.Client.Set(ctx, key, value, expiration)
}

func (r *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.Client.Get(ctx, key)
}

func (r *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.Client.Del(ctx, keys...)
}

func (r *Client) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return r.Client.HSet(ctx, key, values...)
}

func (r *Client) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return r.Client.HGet(ctx, key, field)
}

func (r *Client) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	return r.Client.HGetAll(ctx, key)
}
