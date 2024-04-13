package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Port    string
	Host    string
	DbPg    DBConfigPg
	DbRedis DBConfigRedis
}

type DBConfigPg struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	SslMode  string
}

type DBConfigRedis struct {
	Host string
	Port int
}

func NewConfig() *Config {
	portPg, err := strconv.Atoi(goDotEnvVariable("DB_PG_PORT"))
	if err != nil {
		return nil
	}

	portRedis, err := strconv.Atoi(goDotEnvVariable("DB_REDIS_PORT"))
	if err != nil {
		return nil
	}

	return &Config{
		Port: goDotEnvVariable("PORT"),
		Host: goDotEnvVariable("HOST"),
		DbPg: DBConfigPg{
			User:     goDotEnvVariable("DB_PG_USER"),
			Dbname:   goDotEnvVariable("DB_PG_NAME"),
			Host:     goDotEnvVariable("DB_PG_HOST"),
			SslMode:  "false",
			Password: goDotEnvVariable("DB_PG_PASSWORD"),
			Port:     portPg,
		},
		DbRedis: DBConfigRedis{
			Host: goDotEnvVariable("DB_REDIS_HOST"),
			Port: portRedis,
		},
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load(".env")
	return os.Getenv(key)
}
