package service

import (
	"context"
	"fmt"
	"time"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/logger"

	"github.com/redis/go-redis/v9"
)

type RedisService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl int) error
	Delete(ctx context.Context, key string) error
}

type redisService struct {
	client *redis.Client
	log    logger.Logger
}

func NewRedisService(client *redis.Client, log logger.Logger) *redisService {
	return &redisService{
		client: client,
		log:    log,
	}
}

func (r *redisService) Get(ctx context.Context, key string) (string, error) {
	r.log.Info("Redis GET: key=%s", key)
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf(constant.ErrorDataNotFound)
		}
		r.log.Error("Redis GET error: %v", err)
		return "", fmt.Errorf(constant.ErrorServerGet)
	}
	return result, nil
}

func (r *redisService) Set(ctx context.Context, key string, value string, ttl int) error {
	r.log.Info("Redis SET: key=%s", key)
	err := r.client.Set(ctx, key, value, time.Duration(ttl)*time.Minute).Err()
	if err != nil {
		r.log.Error("Redis SET error: %v", err)
		return fmt.Errorf(constant.ErrorServerCreate)
	}
	return nil
}

func (r *redisService) Delete(ctx context.Context, key string) error {
	r.log.Info("Redis DELETE: key=%s", key)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.log.Error("Redis DELETE error: %v", err)
		return fmt.Errorf(constant.ErrorServerUpdate)
	}
	return nil
}
