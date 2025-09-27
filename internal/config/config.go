package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	Sqlite   SqliteConfig   `yaml:"sqlite"`
	Redis    RedisConfig    `yaml:"redis"`
	SMTP     SMTPConfig     `yaml:"smtp"`
}

type ServerConfig struct {
	Debug bool `yaml:"debug"`
	Port  int  `yaml:"port"`
}

type SecurityConfig struct {
	SessionSecret string `yaml:"session_secret"`
	CSRFSecret    string `yaml:"csrf_secret"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	Debug    bool   `yaml:"debug"`
}

type SqliteConfig struct {
	Path string `yaml:"path"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Sender   string `yaml:"sender"`
}

var global AppConfig

func GetAppConfig() AppConfig {
	return global
}

func InitConfig(path string) AppConfig {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	configBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &global)
	if err != nil {
		panic(err)
	}
	return global
}
