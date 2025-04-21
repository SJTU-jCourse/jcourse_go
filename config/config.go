package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Server Server
	DB     DB
	Redis  Redis
	LLM    LLM
	SMTP   SMTP
	Auth   Auth
}

type DB struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_DBNAME"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
}

type Server struct {
	Debug       bool `env:"DEBUG"`
	NoLoginMode bool `env:"NO_LOGIN_MODE"`
	Port        int  `env:"PORT"`
}

type LLM struct {
	BaseUrl string `env:"LLM_BASE_URL"`
	Token   string `env:"LLM_TOKEN"`
	Model   string `env:"LLM_MODEL"`
}

type SMTP struct {
	Host     string `env:"SMTP_HOST"`
	Port     int    `env:"SMTP_PORT"`
	Username string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
	Sender   string `env:"SMTP_SENDER"`
}

type Auth struct {
	SessionSecret string `env:"SESSION_SECRET"`
	CSRFSecret    string `env:"CSRF_SECRET"`
	HashSalt      string `env:"HASH_SALT"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
