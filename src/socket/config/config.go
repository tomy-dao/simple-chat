package config

import (
	"os"
	"strconv"
)

type ServiceConfig struct {
	HTTPPort int
	Host     string
	

	BackendServerURL string
}

var Config = ServiceConfig{}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func LoadConfig() {
	httpPortStr := getEnv("HTTP_PORT", "8080")
	httpPortInt, err := strconv.Atoi(httpPortStr)
	if err != nil {
		httpPortInt = 80 // fallback default
	}

	host := getEnv("HOST", "0.0.0.0")


	backendServerURL := getEnv("BACKEND_SERVER_URL", "http://localhost")

	Config = ServiceConfig{
		HTTPPort: httpPortInt,
		Host:     host,
		BackendServerURL: backendServerURL,
	}
}
