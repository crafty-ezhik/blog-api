package jwt

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type RedisBlackList struct {
	client *redis.Client
}

type RedisVersioner struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) (*RedisBlackList, *RedisVersioner) {
	return &RedisBlackList{client: client}, &RedisVersioner{client: client}
}

func (r *RedisBlackList) IsBlackListed(token string) (bool, error) {
	val, err := r.client.Get(context.Background(), "jwt_refresh:"+token).Result()
	return val == "revoked", err
}

func (r *RedisBlackList) AddToBlackList(token string, ttl time.Duration) error {
	return r.client.Set(context.Background(), "jwt_refresh:"+token, "revoked", ttl).Err()
}

func (r *RedisVersioner) IncrementVersion(userID uint) error {
	return r.client.Incr(context.Background(), "user_version:"+strconv.Itoa(int(userID))).Err()
}

func (r *RedisVersioner) GetVersion(userID uint) (uint, error) {
	val, err := r.client.Get(context.Background(), "user_version:"+strconv.Itoa(int(userID))).Result()
	if val == "" || err != nil {
		return 0, nil
	}
	version, err := strconv.ParseUint(val, 10, 32)
	return uint(version), err
}
