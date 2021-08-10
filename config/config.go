package config

import "os"

type Config struct {
	VerificationToken string
	WebHookURL string
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
	}
}

func Get() Config {
	return config
}