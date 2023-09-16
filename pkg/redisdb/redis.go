package redisdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type RedisDB struct {
	connAttempts int
	connTimeout  time.Duration
	Client       *redis.Client
}

type Option func(*RedisDB)

func MustRedisDB(url string, opts ...Option) *RedisDB {
	p, err := NewRedisDB(url, opts...)
	if err != nil {
		panic(err)
	}
	return p
}

func NewRedisDB(url string, opts ...Option) (*RedisDB, error) {
	optsRedisClient, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	rdb := &RedisDB{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
		Client:       redis.NewClient(optsRedisClient),
	}

	for _, opt := range opts {
		opt(rdb)
	}

	for rdb.connAttempts > 0 {
		_, err := rdb.Client.Ping(context.Background()).Result()
		if nil == err {
			break
		}

		log.Printf("redis is trying to connect, attempts left: %d", rdb.connAttempts)

		time.Sleep(rdb.connTimeout)

		rdb.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("redis - connAttempts == 0: %w", err)
	}

	return rdb, nil
}

func (r *RedisDB) Close() {
	if r.Client != nil {
		err := r.Client.Close()
		if err != nil {
			log.Printf("redis - error when closing connection: %v", err)
		}
	}
}
