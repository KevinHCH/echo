package redis

import (
	"context"
	utils "echo/internal"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	redisUrl, err := utils.GetEnv("REDIS_URL")
	if err != nil {
		log.Fatalln("failed to load redis URL: %w", err)

	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	return &RedisClient{client: rdb}, nil
}
func (r *RedisClient) Enqueue(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		log.Printf("Error enqueuing key %s: %v", key, err)
		return err
	}
	return nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Error getting key %s: %v", key, err)
		return "", err
	}
	return value, nil
}
func (r *RedisClient) KeyExists(ctx context.Context, key string) (bool, error) {
	// Redis EXISTS returns an integer, 1 if the key exists, 0 if it does not
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		log.Printf("Error checking existence of key %s: %v", key, err)
		return false, err
	}

	return exists == 1, nil
}
func (r *RedisClient) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		log.Printf("Error getting TTL for key %s: %v", key, err)
		return 0, err
	}
	return ttl, nil
}
func (r *RedisClient) GetAll(ctx context.Context) ([]byte, error) {

	var cursor uint64
	var keys []string
	var allKeys []string

	// get all available keys
	for {
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			log.Printf("Error scanning keys: %v", err)
			return nil, err
		}
		allKeys = append(allKeys, keys...)
		if cursor == 0 {
			break
		}
	}

	// Return null if no keys are found
	if len(allKeys) == 0 {
		log.Println("No items found in Redis.")
		return json.Marshal(nil)
	}

	var result []map[string]string

	// Retrieve values for each key
	for _, key := range allKeys {
		// Check the type of the value
		keyType, err := r.client.Type(ctx, key).Result()
		if err != nil {
			log.Printf("Error getting type for key %s: %v", key, err)
			continue
		}

		// Only handle string keys
		if keyType == "string" {
			value, err := r.client.Get(ctx, key).Result()
			if err != nil {
				log.Printf("Error getting value for key %s: %v", key, err)
				continue
			}

			result = append(result, map[string]string{
				"key":   key,
				"value": value,
			})
		} else {
			log.Printf("Skipping non-string key %s of type %s", key, keyType)
		}
	}

	// Convert the result to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return nil, err
	}

	return jsonData, nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
