package jwt

type JWTService struct {
	blackLister BlackListStorage
	versioner   TokenVersionStorage
}

func NewJWTService(bl BlackListStorage, tv TokenVersionStorage) *JWTService {
	return &JWTService{
		blackLister: bl,
		versioner:   tv,
	}
}
