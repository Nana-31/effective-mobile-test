package config

import (
	"log"
	"os"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	ServerPort string
}

func MustLoad() *Config {
	c := &Config{}
	c.DB.Host = os.Getenv("DB_HOST")
	if c.DB.Host == "" {
		c.DB.Host = "localhost"
	}

	c.DB.Port = os.Getenv("DB_PORT")
	c.DB.User = os.Getenv("DB_USER")
	c.DB.Password = os.Getenv("DB_PASSWORD")
	c.DB.Name = os.Getenv("DB_NAME")

	c.Server.ServerPort = os.Getenv("SERVER_PORT")
	if c.Server.ServerPort == "" {
		c.Server.ServerPort = ":8080"
	}

	if c.DB.User == "" || c.DB.Password == "" || c.DB.Name == "" {
		log.Fatal("Config error: critical database environment variables are missing")
	}

	return c
}

func (c *Config) DBDSN() string {
	return "postgres://" + c.DB.User + ":" + c.DB.Password + "@" + c.DB.Host + ":" + c.DB.Port + "/" + c.DB.Name + "?sslmode=disable"
}
