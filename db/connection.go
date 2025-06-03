package db

import (
	"fmt"
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection(conf *config.Config) *gorm.DB {
	logger.Log.Debug("Get connection to database")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		conf.DB.Host,
		conf.DB.Username,
		conf.DB.Password,
		conf.DB.Database,
		conf.DB.Port,
		conf.DB.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
