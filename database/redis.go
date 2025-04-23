package database

import (
	"context"
	"github.com/mfasdfasdf/kit-framework/config"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var RedisCli *RedisClient

type RedisClient struct {
	cli *redis.Client
}

func InitRedis() {
	if RedisCli != nil {
		return
	}
	client := redis.NewClient(&redis.Options{
		Addr:         config.Configuration.Redis.Url + ":" + strconv.Itoa(config.Configuration.Redis.Port),
		Username:     config.Configuration.Redis.Username,
		Password:     config.Configuration.Redis.Password,
		DB:           config.Configuration.Redis.Db,
		PoolSize:     config.Configuration.Redis.PoolSize,
		MaxIdleConns: config.Configuration.Redis.MaxIdleSize,
		MinIdleConns: config.Configuration.Redis.MinIdleSize,
	})
	RedisCli = &RedisClient{cli: client}
}

func (r *RedisClient) SetString(key string, value string, ttl int) {
	RedisCli.cli.Set(context.Background(), key, value, time.Duration(ttl)*time.Second)
}

func (r *RedisClient) GetString(key string) string {
	result, err := RedisCli.cli.Get(context.Background(), key).Result()
	if err != nil {
		return ""
	}
	return result
}

func (r *RedisClient) SetHash(key string, value map[string]any, ttl int) {
	RedisCli.cli.HSet(context.Background(), key, value)
	RedisCli.cli.Expire(context.Background(), key, time.Duration(ttl)*time.Second)
}

func (r *RedisClient) GetHash(key string) map[string]string {
	result, err := RedisCli.cli.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil
	}
	return result
}

func (r *RedisClient) LPush(key string, value []any, ttl int) {
	RedisCli.cli.LPush(context.Background(), key, value...)
	RedisCli.cli.Expire(context.Background(), key, time.Duration(ttl)*time.Second)
}

func (r *RedisClient) LPop(key string) any {
	result, err := RedisCli.cli.LPop(context.Background(), key).Result()
	if err != nil {
		return nil
	}
	return result
}

func (r *RedisClient) RPush(key string, value []any, ttl int) {
	RedisCli.cli.RPush(context.Background(), key, value...)
	RedisCli.cli.Expire(context.Background(), key, time.Duration(ttl)*time.Second)
}

func (r *RedisClient) RPop(key string) any {
	result, err := RedisCli.cli.RPop(context.Background(), key).Result()
	if err != nil {
		return nil
	}
	return result
}

func (r *RedisClient) LLen(key string) int64 {
	result, err := RedisCli.cli.LLen(context.Background(), key).Result()
	if err != nil {
		return 0
	}
	return result
}
