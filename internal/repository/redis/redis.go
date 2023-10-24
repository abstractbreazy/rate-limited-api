package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"

	redis_db "github.com/abstractbreazy/rate-limited-api/pkg/redis"
)

type Redis struct {
	*redis_db.Redis // underlying redis wrapper instance
}

func New(conf *redis_db.Config) (redi *Redis, err error) {
	var rd *redis_db.Redis
	if rd, err = redis_db.New(conf); err != nil {
		return nil, err
	}

	redi = new(Redis)
	redi.Redis = rd
	return
}

// GetResponse object.
type GetResponse struct {
	Value uint64
	TTL   time.Duration
}

// Get the current value/expiration by Subnet key.
func (r *Redis) Get(ctx context.Context, key string) (res *GetResponse, err error) {
	pipeline := r.Client.Pipeline()
	value := pipeline.Get(ctx, key)
	ttl := pipeline.TTL(ctx, key)
	_, err = pipeline.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var ttlDuration time.Duration
	if d, err := ttl.Result(); err != nil || d == -1 || d == -2 {
		ttlDuration = time.Second * 10
		if err := r.Client.Expire(ctx, key, ttlDuration).Err(); err != nil {
			return nil, err
		}
	} else {
		ttlDuration = d
	}

	var total uint64
	total, err = value.Uint64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	return &GetResponse{
		Value: total,
		TTL:   ttlDuration,
	}, nil
}

// Increment Subnet key per request.
func (r *Redis) Incr(ctx context.Context, key string) (value uint64, err error) {
	return r.Client.Incr(ctx, key).Uint64()
}

// Reset Subnet key.
// Used by drop endpoint
func (r *Redis) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
