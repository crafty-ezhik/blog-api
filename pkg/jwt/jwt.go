package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	ErrUnexpectedSigningMethod = "unexpected signing method"
	ErrInvalidToken            = "invalid token"
	ErrTokenNotFound           = "token not found"
	ErrTokenExpired            = "token expired"
)

type JWTData struct {
	UserId  uint  `json:"user_id"`
	Exp     int64 `json:"exp"`
	Version uint  `json:"version"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateToken(signingKey string, userID uint, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     1,
		"version": 0, // TODO:  Сделать получение из Redis
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyToken(tokenString, signingKey string) (*JWTData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(ErrUnexpectedSigningMethod)
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New(ErrInvalidToken)
	}

	userID, ok := claims["user_id"].(uint)
	if !ok {
		return nil, errors.New(ErrInvalidToken)
	}
	exp := int64(claims["exp"].(float64))
	version := uint(claims["version"].(float64))

	return &JWTData{
		UserId:  userID,
		Exp:     exp,
		Version: version,
	}, nil
}

func Refresh(refreshToken, signingKey string, accessTTL, refreshTTL time.Duration) (*Tokens, error) {
	tokenData, err := VerifyToken(refreshToken, signingKey)
	if err != nil {
		return nil, err
	}
	// Надо проверить exp

	// Надо проверить версию

	// Генерация новой пары ключей
	newAccessToken, err := GenerateToken(signingKey, tokenData.UserId, accessTTL)
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := GenerateToken(signingKey, tokenData.UserId, refreshTTL)
	if err != nil {
		return nil, err
	}

	// Добавление старого refresh токена в BlackList

	// Увеличение версии токена

	// Возврат значений
	return &Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
