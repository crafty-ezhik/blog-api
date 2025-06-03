package models

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Title    string `gorm:"size:255" json:"title"`
	Text     string `gorm:"type:text" json:"text"`
	AuthorID uint   `gorm:"index" json:"author_id"` // внешний ключ на User.ID
	//Author    User           `gorm:"foreignKey:AuthorID" json:"author,omitempty"` // загружается через Preload
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
