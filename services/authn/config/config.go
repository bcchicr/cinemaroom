package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Grpc        Grpc
	Postgres    Postgres
	Redis       Redis
	UserService UserService
	Jwt         Jwt
}

type Grpc struct {
	Port string `env:"GRPC_PORT,required"`
}

type Postgres struct {
	Url string `env:"POSTGRES_URL,required"`
}

type Redis struct {
	Addr     string `env:"REDIS_ADDR,required"`
	Password string `env:"REDIS_PASSWORD,required"`
	DB       int    `env:"REDIS_DB,required"`
}

type UserService struct {
	Url string `env:"USER_SERVICE_URL,required"`
}

type Jwt struct {
	AccessSecret       string `env:"ACCESS_TOKEN_SECRET,required"`
	AccessTTLInMinutes int    `env:"ACCESS_TOKEN_TTL_IN_MINUTES,required"`
	RefreshTTLInDays   int    `env:"REFRESH_TOKEN_TTL_IN_DAYS,required"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Println("WARN: failed to load dotenv file. Use environment variables")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
