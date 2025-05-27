package user

import "gorm.io/gorm"

type UserRepository interface {
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (repo *UserRepositoryImpl) FindByID(userId uint) error {
	return repo.db.First(&userId).Error
}
