package config

import (
	"os"
	"strconv"
)

type ServiceConfig struct {
	HTTPPort int

	Host     string
	
	DBHost   string
	DBPort   string
	DBUser   string
	DBPassword string
	DBName     string

	SocketServerURL string
	SocketToken     string
	JwtSecret       string
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
	httpPortStr := getEnv("HTTP_PORT", "80")
	httpPortInt, err := strconv.Atoi(httpPortStr)
	if err != nil {
		httpPortInt = 80 // fallback default
	}

	host := getEnv("HOST", "0.0.0.0")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "your_root_password")
	dbName := getEnv("DB_NAME", "simple_chat")
	jwtSecret := getEnv("JWT_SECRET", "your_jwt_secret")

	socketServerURL := getEnv("SOCKET_SERVER_URL", "http://localhost:8080")
	socketToken := getEnv("SOCKET_TOKEN", "your_socket_token")

	Config = ServiceConfig{
		HTTPPort: httpPortInt,
		Host:     host,
		DBHost:   dbHost,
		DBPort:   dbPort,
		DBUser:   dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		SocketServerURL: socketServerURL,
		SocketToken:     socketToken,
		JwtSecret:       jwtSecret,
	}
}
