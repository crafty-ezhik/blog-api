package user

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(userId uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(userID uint, updateField *models.User) error
	Delete(userID uint) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	logger.Log.Debug("Init user repository")
	return &UserRepositoryImpl{db: db}
}

func (repo *UserRepositoryImpl) FindByID(userId uint) (*models.User, error) {
	var user *models.User
	result := repo.db.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user *models.User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepositoryImpl) Create(user *models.User) error {
	return repo.db.Create(user).Error
}

func (repo *UserRepositoryImpl) Update(userID uint, updateField *models.User) error {
	return repo.db.Model(&models.User{ID: userID}).Updates(updateField).Error
}

func (repo *UserRepositoryImpl) Delete(userID uint) error {
	return repo.db.Delete(&models.User{}, &userID).Error
}
