package jwt

import (
	"github.com/redis/go-redis/v9"
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

func (r *RedisBlackList) IsBlackListed(token string) (bool, error) {}

func (r *RedisBlackList) AddToBlackList(token string, ttl time.Duration) error {}

func (r *RedisVersioner) IncrementVersion(userID uint) error {}

func (r *RedisVersioner) GetVersion(userID uint) (uint, error) {}
