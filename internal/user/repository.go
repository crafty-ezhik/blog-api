package user

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(userId uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(userID uint) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (repo *UserRepositoryImpl) FindByID(userId uint) (*models.User, error) {
	var user *models.User
	result := repo.db.First(user, userId)
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

func (repo *UserRepositoryImpl) Update(user *models.User) error {
	return repo.db.Save(user).Error
}

func (repo *UserRepositoryImpl) Delete(userID uint) error {
	return repo.db.Delete(&userID).Error
}
