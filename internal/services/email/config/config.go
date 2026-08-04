package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ListenAddress string
	SMTP          SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	UserName string
	Password string
	From     string
}

func LoadConfig() (Config, error) {
	listenAddress := strings.TrimSpace(os.Getenv("EMAIL_LISTEN_ADDRESS"))
	if listenAddress == "" {
		return Config{}, fmt.Errorf("EMAIL_LISTEN_ADDRESS is required")
	}
	port, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SMTP_PORT")))
	if err != nil || port <= 0 {
		return Config{}, fmt.Errorf("SMTP_PORT must be a positive integer")
	}
	name := strings.TrimSpace(os.Getenv("NOTEZY_OFFICIAL_NAME"))
	address := strings.TrimSpace(os.Getenv("NOTEZY_OFFICIAL_GMAIL"))
	config := Config{
		ListenAddress: listenAddress,
		SMTP: SMTPConfig{
			Host:     strings.TrimSpace(os.Getenv("SMTP_HOST")),
			Port:     port,
			UserName: address,
			Password: os.Getenv("NOTEZY_OFFICIAL_GOOGLE_APPLICATION_PASSWORD"),
			From:     name + " <" + address + ">",
		},
	}
	if config.SMTP.Host == "" || config.SMTP.UserName == "" || config.SMTP.Password == "" || name == "" {
		return Config{}, fmt.Errorf("SMTP_HOST, NOTEZY_OFFICIAL_NAME, NOTEZY_OFFICIAL_GMAIL, and NOTEZY_OFFICIAL_GOOGLE_APPLICATION_PASSWORD are required")
	}

	return config, nil
}
