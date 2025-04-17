package awsclient

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const RedisClientKey = "redis-client-key"

//go:generate mockery --name WrappedRedisClient
type WrappedRedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	Close() error
	PSubscribe(ctx context.Context, channels ...string) *redis.PubSub
}

type RedisClient struct {
	Client WrappedRedisClient
}

func GetRedisClient(ctx context.Context, addr string) (context.Context, RedisClient) {
	redisClientForAddrKey := RedisClientKey + "-" + addr
	if existingClient, ok := ctx.Value(redisClientForAddrKey).(RedisClient); ok {
		return ctx, existingClient
	}
	wrapper := RedisClient{
		Client: redis.NewClient(&redis.Options{
			MaxRetries:   -1,
			PoolTimeout:  200 * time.Millisecond,
			MinIdleConns: 10,
			Addr:         addr,
		}),
	}
	return context.WithValue(ctx, redisClientForAddrKey, wrapper), wrapper
}

func (wrapper RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	resultCmd := wrapper.Client.Get(ctx, key)
	if resultCmd.Err() != nil {
		return nil, resultCmd.Err()
	}
	return resultCmd.Bytes()
}

func (wrapper RedisClient) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return wrapper.Client.Set(ctx, key, value, expiration).Err()
}
func (wrapper RedisClient) SetNX(ctx context.Context, key string, value []byte, expiration time.Duration) (bool, error) {
	cmd := wrapper.Client.SetNX(ctx, key, value, expiration)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val(), nil
}

func (wrapper RedisClient) Del(ctx context.Context, key string) error {
	return wrapper.Client.Del(ctx, key).Err()
}

func (wrapper RedisClient) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return wrapper.Client.Scan(ctx, cursor, match, count)
}

func (wrapper RedisClient) Close() error {
	return wrapper.Client.Close()
}
