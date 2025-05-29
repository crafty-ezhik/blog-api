package user

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(userId uint) (models.User, error)
	FindByEmail(email string) (models.User, error)
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

func (service *UserRepositoryImpl) FindByID(userId uint) (models.User, error) {
	var user models.User
	result := service.db.First(&user, userId)
	return user, result.Error
}

func (service *UserRepositoryImpl) FindByEmail(email string) (models.User, error) {
	var user models.User
	result := service.db.Where("email = ?", email).First(&user)
	return user, result.Error
}

func (service *UserRepositoryImpl) Create(user *models.User) error {
	return service.db.Create(user).Error
}

func (service *UserRepositoryImpl) Update(user *models.User) error {
	return service.db.Save(user).Error
}

func (service *UserRepositoryImpl) Delete(userID uint) error {
	return service.db.Delete(&userID).Error
}
