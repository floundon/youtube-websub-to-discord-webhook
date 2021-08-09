package config

import "os"

type Config struct {
	VerificationToken string
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
	}
}

func Get() Config {
	return config
}