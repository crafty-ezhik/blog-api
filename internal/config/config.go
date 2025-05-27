package config

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Auth   AuthConfig   `mapstructure:"jwt"`
	DB     DbConfig     `mapstructure:"database"`
	Server ServerConfig `mapstructure:"server"`
}

type AuthConfig struct {
	SigningKey string        `mapstructure:"signing_key"`
	SecretKey  string        `mapstructure:"secret_key"`
	AccessTTL  time.Duration `mapstructure:"access_ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}

type DbConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

func LoadConfig(path string) (*Config, error) {
	// Подгрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// Указание имени и расширения файла
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")

	// Указание путей, где искать конфигурационный файл
	viper.AddConfigPath(path)

	// Загрузка переменных окружения
	viper.SetEnvPrefix("APP")

	// замена .на _в именах ключей, поэтому вложенные ключи, такие как, app.port могут напрямую сопоставляться с переменными среды
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var config *Config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Println("Config file not found, using environment variables only")
		} else {
			log.Fatalf("Ошибка чтения файла конфигурации, %s", err)
		}
	}

	err := viper.Unmarshal(&config, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "mapstructure"
		dc.Result = &config
		dc.WeaklyTypedInput = true // позволяет преобразовывать строки в числа и т.п.
	})
	if err != nil {
		log.Println("Ошибка при анмаршалинге в структуру")
		return nil, err
	}

	return config, nil
}
