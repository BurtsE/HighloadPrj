package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"highload/internal/model"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}
	return &RedisClient{Client: rdb}
}

func SaveMetric(r *RedisClient, m model.Metric) error {
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("error marshaling metric: %v", err)
	}

	// Сохраняем в список (для скользящего окна)
	if err = r.Client.LPush(ctx, "metrics", data).Err(); err != nil {
		return fmt.Errorf("error saving metric: %v", err)
	}
	if err = r.Client.LTrim(ctx, "metrics", 0, 49).Err(); err != nil { // Keep last 50 entries
		return fmt.Errorf("error deleting metric: %v", err)
	}

	return nil
}

func GetMetrics(r *RedisClient) ([]model.Metric, error) {
	values, err := r.Client.LRange(ctx, "metrics", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting metrics: %v", err)

	}

	metrics := make([]model.Metric, len(values))
	for i, v := range values {
		var m model.Metric
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, fmt.Errorf("error unmarshaling metric: %v", err)
		}
		metrics[i] = m
	}
	return metrics, nil
}
