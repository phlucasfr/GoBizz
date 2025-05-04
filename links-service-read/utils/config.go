package utils

import (
	"os"
)

type Config struct {
	FrontendSource string
	DynamoEndpoint string
}

var (
	ConfigInstance Config
)

// LoadEnvInstance initializes the global ConfigInstance with environment variables.
// It retrieves the following environment variables:
// - FRONTEND_SOURCE: The source URL for the frontend.
// - DYNAMODB_ENDPOINT: The endpoint URL for DynamoDB.
// These values are used to populate the Config struct.
func LoadEnvInstance() {
	ConfigInstance = Config{
		FrontendSource: os.Getenv("FRONTEND_SOURCE"),
		DynamoEndpoint: os.Getenv("DYNAMODB_ENDPOINT"),
	}
}
