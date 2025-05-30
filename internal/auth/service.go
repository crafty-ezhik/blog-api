package auth

import (
	"errors"
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/internal/user"
	cjwt "github.com/crafty-ezhik/blog-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrUserExisted        = "the user with this email already exists"
	ErrInvalidCredentials = "invalid credentials"
)

type AuthService interface {
	Register(data *RegisterRequest) (bool, error)
	Login(data *LoginRequest) (*LoginResponse, error)
	Refresh(tokenStr string) (*cjwt.Tokens, error)
	Logout(tokenStr string) error
}

type AuthServiceimpl struct {
	cfg      *config.Config
	jwtAuth  *cjwt.JWT
	UserRepo user.UserRepository
}

func NewAuthService(cfg *config.Config, userRepo user.UserRepository, jwtAuth *cjwt.JWT) *AuthServiceimpl {
	return &AuthServiceimpl{cfg: cfg, UserRepo: userRepo, jwtAuth: jwtAuth}
}

func (s *AuthServiceimpl) Login(data *LoginRequest) (*LoginResponse, error) {
	existedUser, err := s.UserRepo.FindByEmail(data.Email)
	if err != nil || existedUser == nil {
		return nil, errors.New(ErrInvalidCredentials)
	}

	hashedPassword := existedUser.Password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(data.Password))
	if err != nil {
		return nil, errors.New(ErrInvalidCredentials)
	}

	accessToken, err := s.jwtAuth.GenerateToken(existedUser.ID, cjwt.Access)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtAuth.GenerateToken(existedUser.ID, cjwt.Refresh)
	if err != nil {
		return nil, err
	}

	output := &LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}

	return output, nil
}

func (s *AuthServiceimpl) Register(data *RegisterRequest) (bool, error) {
	existedUser, _ := s.UserRepo.FindByEmail(data.Email)
	if existedUser != nil {
		return false, errors.New(ErrUserExisted)
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	newUser := &models.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: string(hashedPass),
		Age:      data.Age,
	}
	err = s.UserRepo.Create(newUser)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *AuthServiceimpl) Refresh(tokenStr string) (*cjwt.Tokens, error) {
	return s.jwtAuth.Refresh(tokenStr)
}

func (s *AuthServiceimpl) Logout(tokenStr string) error {
	return s.jwtAuth.Logout(tokenStr)
}
