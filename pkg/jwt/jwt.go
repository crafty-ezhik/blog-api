package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenType int

const (
	Access TokenType = iota
	Refresh
)

const (
	ErrUnexpectedSigningMethod = "unexpected signing method"
	ErrInvalidToken            = "invalid token"
	ErrTokenNotFound           = "token not found"
	ErrSessionExpired          = "session expired"
	ErrInternalServer          = "internal server error"
	ErrRefreshExpired          = "refresh token expired due to logout / password change"
	ErrUnknownTokenType        = "unknown token type"
	ErrInBlackList             = "refresh token revoked or not found"
)

type JWT struct {
	jwtService *JWTService
	accessTTL  time.Duration
	refreshTTL time.Duration
	signingKey string
}

func NewJWT(jwtService *JWTService, accessTTL, refreshTTL time.Duration, signingKey string) *JWT {
	return &JWT{
		jwtService: jwtService,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		signingKey: signingKey,
	}
}

type JWTData struct {
	UserId  uint  `json:"user_id"`
	Exp     int64 `json:"exp"`
	Version uint  `json:"version"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (j *JWT) GenerateToken(userID uint, tokenType TokenType) (string, error) {
	currentVersion, err := j.jwtService.versioner.GetVersion(userID)
	if err != nil {
		// TODO: Добавить логирование ошибки
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now().Unix(),
		"version": currentVersion,
	}

	switch tokenType {
	case Access:
		claims["exp"] = time.Now().Add(j.accessTTL).Unix()
	case Refresh:
		claims["exp"] = time.Now().Add(j.refreshTTL).Unix()
	default:
		return "", errors.New(ErrUnknownTokenType)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.signingKey))
	if err != nil {
		// TODO: Добавить логирование ошибки
		return "", err
	}
	return signedToken, nil
}

func (j *JWT) VerifyToken(tokenString string) (*JWTData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(ErrUnexpectedSigningMethod)
		}
		return []byte(j.signingKey), nil
	})
	if err != nil {
		// TODO: Добавить логирование ошибки
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		// TODO: Добавить логирование ошибки
		return nil, errors.New(ErrInvalidToken)
	}

	// Получение UserID
	userID, ok := claims["user_id"].(uint)
	if !ok {
		// TODO: Добавить логирование ошибки
		return nil, errors.New(ErrInvalidToken)
	}

	// Получение и проверка версии
	version := uint(claims["version"].(float64))
	currentVersion, err := j.jwtService.versioner.GetVersion(userID)
	if err != nil {
		// TODO: Добавить логирование ошибки
		return nil, errors.New(ErrInternalServer)
	}
	if version < currentVersion {
		// TODO: Добавить логирование ошибки
		return nil, errors.New(ErrRefreshExpired)
	}

	// Получение и проверка exp
	exp := int64(claims["exp"].(float64))
	if time.Now().Unix() > exp {
		// TODO: Добавить логирование ошибки
		return nil, errors.New(ErrSessionExpired)
	}

	return &JWTData{
		UserId:  userID,
		Exp:     exp,
		Version: version,
	}, nil
}

func (j *JWT) Refresh(refreshToken string) (*Tokens, error) {
	// Надо проверить не черном ли списке токен
	if ok, err := j.jwtService.blackLister.IsBlackListed(refreshToken); err != nil || !ok {
		// TODO: Добавить логирование ошибки
		return nil, errors.New(ErrInBlackList)
	}

	// Парсинг токена
	tokenData, err := j.VerifyToken(refreshToken)
	if err != nil {
		// TODO: Добавить логирование ошибки
		return nil, err
	}

	// Генерация новой пары ключей
	newAccessToken, err := j.GenerateToken(tokenData.UserId, Access)
	if err != nil {
		// TODO: Добавить логирование ошибки
		return nil, err
	}
	newRefreshToken, err := j.GenerateToken(tokenData.UserId, Refresh)
	if err != nil {
		// TODO: Добавить логирование ошибки
		return nil, err
	}

	// Добавление старого refresh токена в BlackList
	err = j.jwtService.blackLister.AddToBlackList(refreshToken, time.Until(time.Unix(tokenData.Exp, 0)))
	if err != nil {
		// TODO: Добавить логирование ошибки
		return nil, err
	}

	// Увеличение версии токена
	err = j.jwtService.versioner.IncrementVersion(tokenData.UserId)
	if err != nil {
		// TODO: Добавить логирование ошибки
		return nil, err
	}

	// Возврат значений
	return &Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (j *JWT) Logout(refreshToken string) error {
	if ok, err := j.jwtService.blackLister.IsBlackListed(refreshToken); err != nil || !ok {
		// TODO: Добавить логирование ошибки
		return errors.New(ErrInBlackList)
	}

	// Парсинг токена
	tokenData, err := j.VerifyToken(refreshToken)
	if err != nil {
		// TODO: Добавить логирование ошибки
		return err
	}

	err = j.jwtService.versioner.IncrementVersion(tokenData.UserId)
	if err != nil {
		// TODO: Добавить логирование ошибки
		return errors.New(ErrInternalServer)
	}

	err = j.jwtService.blackLister.AddToBlackList(refreshToken, time.Until(time.Unix(tokenData.Exp, 0)))
	if err != nil {
		// TODO: Добавить логирование ошибки
		return errors.New(ErrInternalServer)
	}
	return nil
}
