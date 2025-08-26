package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var cfg *Config

type Config struct {
	EMQX_host     string
	EMQX_port     string
	EMQX_username string
	EMQX_password string

	ApiKey string

	SensorType string
	ID1        string
	ID2        int
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	id2, _ := strconv.Atoi(os.Getenv("SENSOR_ID2"))

	cfg := &Config{
		EMQX_host:     os.Getenv("EMQX_HOST"),
		EMQX_port:     os.Getenv("EMQX_PORT"),
		EMQX_username: os.Getenv("EMQX_USERNAME"),
		EMQX_password: os.Getenv("EMQX_PASSWORD"),

		ApiKey: os.Getenv("API_KEY"),

		SensorType: os.Getenv("SENSOR_TYPE"),
		ID1:        os.Getenv("SENSOR_ID1"),
		ID2:        id2,
	}
	return cfg, nil
}
