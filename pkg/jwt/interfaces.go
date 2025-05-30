package jwt

import "time"

type BlackListStorage interface {
	IsBlackListed(token string) (bool, error)
	AddToBlackList(token string, ttl time.Duration) error
}

type TokenVersionStorage interface {
	IncrementVersion(userID uint) error
	GetVersion(userID uint) (uint, error)
}
