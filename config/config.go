package config

import "os"

type Config struct {
	VerificationToken         string
	WebHookURL                string
	YouTubeVideoDataTableName string
}

var config Config

func init() {
	config = Config{
		VerificationToken: func() string {
			envToken := os.Getenv("WEBSUB_VERIFICATION_TOKEN")
			if envToken == "" {
				panic("verification token is not specified")
			}
			return envToken
		}(),

		WebHookURL: func() string {
			envURL := os.Getenv("WEBHOOK_URL")
			if envURL == "" {
				panic("webhook url is not specified")
			}
			return envURL
		}(),

		YouTubeVideoDataTableName: func() string {
			tableName := os.Getenv("YOUTUBE_VIDEO_DATA_TABLE_NAME")
			if tableName == "" {
				panic("youtube video data table name is not specified")
			}
			return tableName
		}(),
	}
}

func Get() Config {
	return config
}
