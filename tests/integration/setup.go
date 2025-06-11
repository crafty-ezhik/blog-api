package integration

import (
	"fmt"
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const TestDatabaseName = "blog_api_test"

func SetupTestDB() *gorm.DB {
	cfg, err := config.LoadConfig("C:\\Users\\Alexandr\\GolandProjects\\blog-api\\tests\\integration\\configs")
	if err != nil {
		log.Errorf("Error loading configs: %v", err)
	}

	cfg.DB.Database = TestDatabaseName

	return getConnectionDB(cfg)
}

func TeardownTestDB(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		log.Errorf("Error dropping table: %v", err)
	}
}

func MigrateTables(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Comment{}, &models.Post{})
	if err != nil {
		panic(err)
	}
}

func CleanupTables(db *gorm.DB) {
	fmt.Println("Cleaning up test DB")
}

func getConnectionDB(conf *config.Config) *gorm.DB {
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
