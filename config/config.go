package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Telegram `yaml:"telegram"`
		Yandex   `yaml:"yandex"`
		Env      `yaml:"env"`
	}

	Telegram struct {
		AllowChats    	[]string 	`env-required:"true" yaml:"allow_chats" env:"ALLOW_CHATS"`
		AllowUsers    	[]string 	`env-required:"true" yaml:"allow_users" env:"ALLOW_USERS"`
		AlertChat		string		`env-required:"true" yaml:"alert_chat" env:"ALERT_CHAT"`
		TelegramToken 	string   	`env:"TELEGRAM_TOKEN"`
	}

	Yandex struct {
		YandexOauth string `env:"YANDEX_OAUTH"`
	}

	Env struct {
		Environment string `env-required:"true" yaml:"environment" env:"ENVIRONMENT" env-default:"prod"`
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	err := cleanenv.ReadConfig(path, config)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
