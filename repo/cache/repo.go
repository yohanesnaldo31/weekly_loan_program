package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisProvider interface {
	// Del Redis `Del key` command. It returns 1 when key deleted.
	Del(ctx context.Context, keys ...string) *redis.IntCmd

	// Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
	Get(ctx context.Context, key string) *redis.StringCmd

	// Set Redis `SET key value [expiration]` command.
	// Use expiration for `SETEx`-like behavior.
	//
	// Zero expiration means the key has no expiration time.
	// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
	// otherwise you will receive an error: (error) ERR syntax error.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type Repository struct {
	redis redisProvider
}

func NewRepository(redis redisProvider) *Repository {
	return &Repository{
		redis: redis,
	}
}
