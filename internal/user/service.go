package user

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
)

type UserService interface {
	GetByID(userID uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(userID uint, updatedFields *models.User) error
	Delete(userID uint) error
}

type UserServiceImpl struct {
	UserRepo UserRepository
}

func NewUserService(UserRepo UserRepository) *UserServiceImpl {
	logger.Log.Debug("Init user service")
	return &UserServiceImpl{UserRepo: UserRepo}
}

func (s *UserServiceImpl) GetByID(userID uint) (*models.User, error) {
	return s.UserRepo.FindByID(userID)
}

func (s *UserServiceImpl) GetByEmail(email string) (*models.User, error) {
	return s.UserRepo.FindByEmail(email)
}

func (s *UserServiceImpl) Create(user *models.User) error {
	return s.UserRepo.Create(user)
}

func (s *UserServiceImpl) Update(userID uint, updatedFields *models.User) error {
	return s.UserRepo.Update(userID, updatedFields)
}

func (s *UserServiceImpl) Delete(userID uint) error {
	return s.UserRepo.Delete(userID)
}
