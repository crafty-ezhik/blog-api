package auth

import (
	"errors"
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(data *RegisterRequest) (bool, error)
	Login(email string, password string) (string, error)
}

type AuthServiceimpl struct {
	cfg      *config.Config
	UserRepo user.UserRepository
}

func NewAuthService(cfg *config.Config, userRepo user.UserRepository) *AuthServiceimpl {
	return &AuthServiceimpl{cfg: cfg, UserRepo: userRepo}
}

func (s *AuthServiceimpl) Login(email string, password string) (string, error) {
	return "", nil
}

func (s *AuthServiceimpl) Register(data *RegisterRequest) (bool, error) {
	existedUser, _ := s.UserRepo.FindByEmail(data.Email)
	if existedUser != nil {
		return false, errors.New("the user with this email already exists")
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
