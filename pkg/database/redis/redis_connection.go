package redis

import (
    "github.com/go-redis/redis/v8"
    "os"
)

var RedisClient *redis.Client

func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),   
        Password: "",
        DB:       0,                            // Default DB
    })
}
