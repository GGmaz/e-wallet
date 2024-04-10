package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Port     string
	Host     string
	Db       DBConfig
	NatsUrl  string
	KafkaUrl string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	SslMode  string
}

func NewConfig() *Config {
	port, err := strconv.Atoi(goDotEnvVariable("DB_PORT"))

	if err != nil {
		return nil
	}

	return &Config{
		Port: goDotEnvVariable("PORT"),
		Host: goDotEnvVariable("HOST"),
		Db: DBConfig{
			User:     goDotEnvVariable("DB_USER"),
			Dbname:   goDotEnvVariable("DB_NAME"),
			Host:     goDotEnvVariable("DB_HOST"),
			SslMode:  "false",
			Password: goDotEnvVariable("DB_PASSWORD"),
			Port:     port,
		},
		//NatsUrl:  goDotEnvVariable("NATS_URL"),
		//KafkaUrl: goDotEnvVariable("KAFKA_URL"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load(".env")
	return os.Getenv(key)
}
