package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"unique" json:"email"`
	Password  string         `json:"-"`
	Age       int            `json:"age,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Один ко многим
	Posts    []Post    `gorm:"foreignKey:AuthorID" json:"posts"`    // один пользователь - много постов
	Comments []Comment `gorm:"foreignKey:AuthorID" json:"comments"` // один пользователь - много комментариев
}
