package redis

import (
	"context"
	"os"
	"testing"
	"time"

	redis_db "github.com/abstractbreazy/rate-limited-api/pkg/redis"
	"github.com/stretchr/testify/require"
)

func connectRedis(t *testing.T) *Redis {
	var url, found = os.LookupEnv("TEST_REDIS_URL")
	if !found {
		url = "redis://localhost:6379"
	}

	rd, err := New(&redis_db.Config{URL: url})
	require.NoError(t, err)
	return rd
}

var (
	testKey          = "123.45.67.0/24"
	testValue uint64 = 1
	testTTL          = 10 * time.Second
)

func TestRedisGet(t *testing.T) {
	// Create a Redis instance
	redis := connectRedis(t)
	defer redis.Close()

	ctx := context.Background()
	err := redis.Client.Set(ctx, testKey, testValue, testTTL).Err()
	require.NoError(t, err)

	res, err := redis.Get(ctx, testKey)
	require.NoError(t, err)
	require.Equal(t, testValue, res.Value)
	require.Equal(t, testTTL, res.TTL)
}

func TestRedisIncr(t *testing.T) {
	redis := connectRedis(t)
	defer redis.Close()

	ctx := context.Background()
	res, err := redis.Incr(ctx, testKey)
	require.NoError(t, err)
	require.Equal(t, uint64(2), res)
}
