package config

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	Host        string `envconfig:"DB_HOST" required:"true"`
	Port        string `envconfig:"DB_PORT" required:"true"`
	User        string `envconfig:"DB_USER" required:"true"`
	Password    string `envconfig:"DB_PASSWORD" required:"true"`
	Name        string `envconfig:"DB_NAME" required:"true"`
	LogMode     bool   `envconfig:"DB_LOG_MODE" required:"true"`
	MaxIdle     int    `envconfig:"DB_MAX_IDLE_CONNS" required:"true"`
	MaxOpen     int    `envconfig:"DB_MAX_OPEN_CONNS" required:"true"`
	MaxLife     int    `envconfig:"DB_MAX_LIFE_TIME" required:"true"`
	MaxIdleTime int    `envconfig:"DB_MAX_IDLE_TIME" required:"true"`
}

type TokenConfig struct {
	IssuerName              string                 `envconfig:"TOKEN_ISSUE" required:"true"`
	JwtSignatureKy          []byte                 `envconfig:"TOKEN_SECRET" required:"true"`
	JwtExpiresTime          time.Duration          `envconfig:"TOKEN_EXPIRE" required:"true"`
	RefreshTokenExpiresTime time.Duration          `envconfig:"REFRESH_TOKEN_EXPIRE" required:"true"`
	JwtSigningMethod        *jwt.SigningMethodHMAC `ignored:"true"`
}

type PortConfig struct {
	ServerPort string `envconfig:"SERVER_PORT" required:"true"`
}

type Config struct {
	DBConfig
	TokenConfig
	PortConfig
}

func (c *Config) validate() error {
	if len(c.TokenConfig.JwtSignatureKy) == 0 {
		return fmt.Errorf("JWT signature key is required")
	}
	if c.TokenConfig.JwtExpiresTime <= 0 {
		return fmt.Errorf("JWT expiration time must be greater than zero")
	}

	c.TokenConfig = TokenConfig{
		IssuerName:              c.TokenConfig.IssuerName,
		JwtSignatureKy:          c.TokenConfig.JwtSignatureKy,
		JwtExpiresTime:          c.TokenConfig.JwtExpiresTime,
		RefreshTokenExpiresTime: c.TokenConfig.RefreshTokenExpiresTime,
		JwtSigningMethod:        c.TokenConfig.JwtSigningMethod,
	}
	return nil
}

func NewConfig() *Config {
	var cfg Config
	envconfig.MustProcess("", &cfg)
	if err := cfg.validate(); err != nil {
		panic(err)
	}
	cfg.TokenConfig.JwtSigningMethod = jwt.SigningMethodHS256

	return &cfg
}
