package main

import (
	"fmt"
	"github.com/crafty-ezhik/blog-api/db"
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/internal/models"
)

func main() {
	fmt.Println("Запуск автоматической миграции БД")
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		fmt.Println(err)
		return
	}
	DB := db.GetConnection(cfg)

	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Миграция успешно проведена!")
}
