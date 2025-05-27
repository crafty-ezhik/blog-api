package models

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Title   string `gorm:"size:255" json:"title"`
	Content string `gorm:"type:text" json:"content"`

	AuthorID uint `gorm:"index" json:"-"`                    // внешний ключ на User.ID
	Author   User `gorm:"foreignKey:AuthorID" json:"author"` // Для получения автора через Preload

	PostID uint `gorm:"index" json:"-"`                // внешний ключ на Post.ID
	Post   Post `gorm:"foreignKey:PostID" json:"post"` // Для получения поста через Preload

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
