package redis

import (
	"time"

	"github.com/go-redis/redis/v7"
)

type Option func(*redis.Options)

func WithAddress(addr string) Option {
	return func(r *redis.Options) {
		r.Addr = addr
	}
}

func WithPassword(pw string) Option {
	return func(r *redis.Options) {
		r.Password = pw
	}
}

func WithDB(db int) Option {
	return func(r *redis.Options) {
		r.DB = db
	}
}

func WithDBTimeout(t time.Duration) Option {
	return func(r *redis.Options) {
		r.DialTimeout = t
	}
}

func New(opts ...Option) (*redis.Client, error) {
	opt := &redis.Options{
		Addr:        "localhost:6379",
		Password:    "",
		DB:          0,
		DialTimeout: 60 * time.Second,
	}
	for _, o := range opts {
		o(opt)
	}

	client := redis.NewClient(opt)

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
