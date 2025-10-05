package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"

	repo "orders/internal/repository"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redisClient *redis.Client
	Capacity    int32
}

const cacheCapacity int32 = 200

func NewCache() *Cache {
	rURL := os.Getenv("REDIS_CONN_STRING")

	opt, err := redis.ParseURL(rURL)
	if err != nil {
		log.Fatalln("Error connecting to redis:", err)
	}

	rdb := redis.NewClient(opt)
	return &Cache{redisClient: rdb, Capacity: cacheCapacity}
}

func (c *Cache) LoadInitialOrders(ctx context.Context, repo *repo.Repository, limit int32) error {
	latestOrders, err := repo.GetLatestOrders(ctx, limit)
	if err != nil {
		return err
	}

	for _, order := range latestOrders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Println("Error marshalling order to JSON:", err)
			continue
		}

		redisKey := order.OrderUID
		err = c.redisClient.Set(ctx, redisKey, orderJSON, 0).Err()
		if err != nil {
			log.Println("Error setting key in Redis:", err)
			continue
		}
	}

	return nil
}
