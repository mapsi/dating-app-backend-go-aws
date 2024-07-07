package config

import "os"

type Config struct {
	JwtSecret      string
	Port           string
	AWSEndpoint    string
	AWSRegion      string
	AWSAccessKeyID string
	AWSSecretKey   string
}

func Load() (*Config, error) {
	return &Config{
		JwtSecret:      getEnv("JWT_SECRET", "super_secret_key"),
		Port:           getEnv("PORT", "3000"),
		AWSEndpoint:    getEnv("AWS_ENDPOINT", "http://localhost:4566"),
		AWSRegion:      getEnv("AWS_REGION", "eu-west-2"),
		AWSAccessKeyID: getEnv("AWS_ACCESS_KEY_ID", "awsAccessKeyId"),
		AWSSecretKey:   getEnv("AWS_SECRET_KEY", "awsSecretKey"),
	}, nil
}

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
