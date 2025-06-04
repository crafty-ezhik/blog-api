package jwt

import "time"

//go:generate mockgen -source=interfaces.go -destination=mock/interfaces_mock.go

type BlackListStorage interface {
	IsBlackListed(token string) bool
	AddToBlackList(token string, ttl time.Duration) error
}

type TokenVersionStorage interface {
	IncrementVersion(userID uint) error
	GetVersion(userID uint) (uint, error)
}
