package config

import (
	"os"
	"strconv"
	"strings"
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

	// Temporal
	TemporalAddress string

	// Kafka
	KafkaBrokers        []string
	KafkaConsumerGroup  string
	KafkaMessageTopic   string
	KafkaNotificationTopic string

	// Rate Limiting
	RateLimitEnabled        bool
	RateLimitRequestsPerMin int
	RateLimitBurst          int
}

var Config = ServiceConfig{}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func splitAndTrim(s string, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
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

	// Temporal configuration
	temporalAddress := getEnv("TEMPORAL_ADDRESS", "localhost:7233")

	// Kafka configuration
	kafkaBrokersStr := getEnv("KAFKA_BROKERS", "localhost:9092")
	kafkaBrokers := []string{kafkaBrokersStr}
	if kafkaBrokersStr != "" {
		// Support multiple brokers separated by comma
		kafkaBrokers = splitAndTrim(kafkaBrokersStr, ",")
	}
	kafkaConsumerGroup := getEnv("KAFKA_CONSUMER_GROUP", "simple-chat-consumer-group")
	kafkaMessageTopic := getEnv("KAFKA_MESSAGE_TOPIC", "chat-messages")
	kafkaNotificationTopic := getEnv("KAFKA_NOTIFICATION_TOPIC", "chat-notifications")

	// Rate limiting configuration
	rateLimitEnabled := getEnv("RATE_LIMIT_ENABLED", "true") == "true"
	rateLimitRequestsPerMin := 60 // default
	if rpm := getEnv("RATE_LIMIT_REQUESTS_PER_MIN", "60"); rpm != "" {
		if val, err := strconv.Atoi(rpm); err == nil {
			rateLimitRequestsPerMin = val
		}
	}
	rateLimitBurst := 10 // default
	if burst := getEnv("RATE_LIMIT_BURST", "10"); burst != "" {
		if val, err := strconv.Atoi(burst); err == nil {
			rateLimitBurst = val
		}
	}

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
		TemporalAddress: temporalAddress,
		KafkaBrokers:        kafkaBrokers,
		KafkaConsumerGroup:  kafkaConsumerGroup,
		KafkaMessageTopic:   kafkaMessageTopic,
		KafkaNotificationTopic: kafkaNotificationTopic,
		RateLimitEnabled:        rateLimitEnabled,
		RateLimitRequestsPerMin: rateLimitRequestsPerMin,
		RateLimitBurst:          rateLimitBurst,
	}
}
