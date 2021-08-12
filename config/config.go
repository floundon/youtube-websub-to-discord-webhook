package config

import (
	"fmt"
	"os"
)

type Config struct {
	VerificationToken         string
	WebHookURL                string
	YouTubeVideoDataTableName string
	YouTubeAPIKey             string
}

var config Config

func init() {
	config = Config{
		VerificationToken:         mustEnvironment("WEBSUB_VERIFICATION_TOKEN"),
		WebHookURL:                mustEnvironment("WEBHOOK_URL"),
		YouTubeVideoDataTableName: mustEnvironment("YOUTUBE_VIDEO_DATA_TABLE_NAME"),
		YouTubeAPIKey:             mustEnvironment("YOUTUBE_API_KEY"),
	}
}

func Get() Config {
	return config
}

func mustEnvironment(env string) string {
	variable := os.Getenv(env)
	if variable == "" {
		panic(fmt.Sprintf("environment variable %s is required", env))
	}
	return variable
}
