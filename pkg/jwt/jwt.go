package jwt

import (
	"errors"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

//go:generate mockgen -source=jwt.go -destination=mock/jwt_mock.go

type JWTInterface interface {
	GenerateToken(userID uint, tokenType TokenType) (string, error)
	VerifyToken(tokenString string) (*JWTData, error)
	Refresh(refreshToken string) (*Tokens, error)
	Logout(refreshToken string) error
}

type TokenType int

const (
	Access TokenType = iota
	Refresh
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrSessionExpired          = errors.New("session expired")
	ErrInternalServer          = errors.New("internal server error")
	ErrRefreshExpired          = errors.New("refresh token expired due to logout / password change")
	ErrUnknownTokenType        = errors.New("unknown token type")
	ErrInBlackList             = errors.New("refresh token revoked or not found")
)

type JWT struct {
	jwtService *JWTService
	accessTTL  time.Duration
	refreshTTL time.Duration
	signingKey string
}

func NewJWT(jwtService *JWTService, accessTTL, refreshTTL time.Duration, signingKey string) *JWT {
	logger.Log.Debug("Init JWT module")
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
	logger.Log.Info("Calling the GenerateToken function")
	logger.Log.Debug("Generating new token", zap.Uint("user_id", userID))
	currentVersion, err := j.jwtService.versioner.GetVersion(userID)
	if err != nil {
		logger.Log.Error("Error generating token: ", zap.Error(err))
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now().Unix(),
		"version": currentVersion,
	}

	logger.Log.Debug("Check token type")
	switch tokenType {
	case Access:
		claims["exp"] = time.Now().Add(j.accessTTL).Unix()
	case Refresh:
		claims["exp"] = time.Now().Add(j.refreshTTL).Unix()
	default:
		return "", ErrUnknownTokenType
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.signingKey))
	if err != nil {
		logger.Log.Error("Error when signing the token: ", zap.Error(err))
		return "", err
	}
	logger.Log.Info("Generated new token successfully")
	return signedToken, nil
}

func (j *JWT) VerifyToken(tokenString string) (*JWTData, error) {
	logger.Log.Info("Calling the VerifyToken function")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(j.signingKey), nil
	})
	if err != nil {
		logger.Log.Error("Error parsing token", zap.Error(err))
		return nil, ErrInvalidToken
	}

	logger.Log.Debug("Conversion to jwt.MapClaims")
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logger.Log.Error("Invalid token")
		return nil, ErrInvalidToken
	}

	// Получение UserID
	logger.Log.Debug("Get user_id")
	userID := uint(claims["user_id"].(float64))

	// Получение и проверка версии
	logger.Log.Debug("Get and check version")
	version := uint(claims["version"].(float64))
	currentVersion, err := j.jwtService.versioner.GetVersion(userID)
	if err != nil {
		logger.Log.Error("Error generating token", zap.Error(err))
		return nil, ErrInternalServer
	}
	if version < currentVersion {
		logger.Log.Debug("The token version does not match:", zap.Error(err))
		return nil, ErrRefreshExpired
	}

	// Получение и проверка exp
	logger.Log.Debug("Get and check exp")
	exp := int64(claims["exp"].(float64))
	if time.Now().Unix() > exp {
		logger.Log.Debug("The token expired")
		return nil, ErrSessionExpired
	}

	return &JWTData{
		UserId:  userID,
		Exp:     exp,
		Version: version,
	}, nil
}

func (j *JWT) Refresh(refreshToken string) (*Tokens, error) {
	logger.Log.Info("Calling the Refresh function")
	// Надо проверить не черном ли списке токен
	if j.jwtService.blackLister.IsBlackListed(refreshToken) {
		logger.Log.Debug("Token is blacklisted")
		return nil, ErrInBlackList
	}

	// Парсинг токена
	tokenData, err := j.VerifyToken(refreshToken)
	if err != nil {
		logger.Log.Error("Error verifying refresh token", zap.Error(err))
		return nil, err
	}

	// Увеличение версии токена
	logger.Log.Debug("Increment token version")
	err = j.jwtService.versioner.IncrementVersion(tokenData.UserId)
	if err != nil {
		logger.Log.Error("Error incrementing version", zap.Error(err))
		return nil, err
	}

	// Генерация новой пары ключей
	newAccessToken, err := j.GenerateToken(tokenData.UserId, Access)
	if err != nil {
		logger.Log.Error("Error generating new access token", zap.Error(err))
		return nil, err
	}
	newRefreshToken, err := j.GenerateToken(tokenData.UserId, Refresh)
	if err != nil {
		logger.Log.Error("Error generating new refresh token", zap.Error(err))
		return nil, err
	}

	// Добавление старого refresh токена в BlackList
	logger.Log.Debug("Add old token into blacklist")
	err = j.jwtService.blackLister.AddToBlackList(refreshToken, time.Until(time.Unix(tokenData.Exp, 0)))
	if err != nil {
		logger.Log.Error("Error adding token to blacklist", zap.Error(err))
		return nil, err
	}

	// Возврат значений
	return &Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (j *JWT) Logout(refreshToken string) error {
	logger.Log.Info("Calling the Logout function")
	if j.jwtService.blackLister.IsBlackListed(refreshToken) {
		logger.Log.Debug("Token is blacklisted")
		return ErrInBlackList
	}

	// Парсинг токена
	tokenData, err := j.VerifyToken(refreshToken)
	if err != nil {
		logger.Log.Error("Error verifying refresh token", zap.Error(err))
		return err
	}

	logger.Log.Debug("Increment token version")
	err = j.jwtService.versioner.IncrementVersion(tokenData.UserId)
	if err != nil {
		logger.Log.Error("Error incrementing version", zap.Error(err))
		return ErrInternalServer
	}

	logger.Log.Debug("Add old token into blacklist")
	err = j.jwtService.blackLister.AddToBlackList(refreshToken, time.Until(time.Unix(tokenData.Exp, 0)))
	if err != nil {
		logger.Log.Error("Error adding token to blacklist", zap.Error(err))
		return ErrInternalServer
	}
	return nil
}
