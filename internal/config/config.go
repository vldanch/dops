package config

type Config struct {
	Ping struct {
		Timeout int
		Retries int
	}
	Notify struct {
		TelegramToken string
		ChatID        string
	}
}

var appConfig Config

func Set(c Config) {
	appConfig = c
}

func Get() Config {
	return appConfig
}
