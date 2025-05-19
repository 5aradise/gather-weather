package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		Env           string `envconfig:"ENV" required:"true"`
		WeatherApiKey string `envconfig:"WEATHERAPI_KEY" required:"true"`
		DB            DB
		Server        Server
		Mail          Mail
	}

	DB struct {
		Address  string `envconfig:"DB_HOST" required:"true"`
		Port     string `envconfig:"DB_PORT" required:"true"`
		Name     string `envconfig:"DB_NAME" required:"true"`
		User     string `envconfig:"DB_USER" required:"true"`
		Password string `envconfig:"DB_PASSWORD" required:"true"`
	}

	Server struct {
		Port string `envconfig:"SERVER_PORT" required:"true"`
	}

	Mail struct {
		Host     string `envconfig:"MAIL_HOST" required:"true"`
		Port     string `envconfig:"MAIL_PORT" required:"true"`
		Sender   string `envconfig:"MAIL_SENDER" required:"true"`
		Password string `envconfig:"MAIL_PASSWORD" required:"true"`
	}
)

func New() (Config, error) {
	var c Config

	err := envconfig.Process("", &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func Load(envPath ...string) error {
	return godotenv.Load(envPath...)
}
