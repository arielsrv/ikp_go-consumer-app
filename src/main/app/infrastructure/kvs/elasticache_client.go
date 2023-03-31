package kvs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/src/main/app/config/env"
)

var ctx = context.Background()

type ElasticCacheClient[TValue any] struct {
	client *redis.Client
}

func NewElasticCacheClient[TValue any](redisClient *redis.Client) *ElasticCacheClient[TValue] {
	return &ElasticCacheClient[TValue]{
		client: redisClient,
	}
}

func (e ElasticCacheClient[TValue]) Get(key string) (*TValue, error) {
	result, err := e.client.
		Get(ctx, key).
		Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	value := new(TValue)
	err = json.Unmarshal([]byte(result), value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (e ElasticCacheClient[TValue]) Save(key string, value *TValue) error {
	if env.IsEmpty(key) {
		return errors.New("missing key")
	}

	if env.IsNil(value) {
		return fmt.Errorf("missing value for key %s", key)
	}

	err := e.client.
		Set(ctx, key, value, 0).
		Err()

	if err != nil {
		return err
	}

	return nil
}
