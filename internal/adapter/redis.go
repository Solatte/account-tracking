package adapter

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	clients = make(map[int]*redis.Client)
	once    sync.Once
)

func InitRedisClients(addr string, password string) error {
	if addr == "" {
		return errors.New("Redis host is empty")
	}

	var initError error
	once.Do(func() {
		for db := 1; db <= 7; db++ {
			client := redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: password,
				DB:       db,
			})

			// Ping the Redis server to check the connection
			if _, err := client.Ping(context.Background()).Result(); err != nil {
				initError = fmt.Errorf("failed to connect to Redis DB %d: %v", db, err)
				return
			}

			clients[db] = client
		}
	})

	return initError
}

func GetRedisClient(db int) (*redis.Client, error) {
	client, exists := clients[db]
	if !exists {
		return nil, fmt.Errorf("redis client for DB %d is not initialized. call InitRedisClients first", db)
	}
	return client, nil
}
