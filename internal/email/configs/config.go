package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"

	emailcontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
)

type Config struct {
	ListenAddress string
	SMTP          SMTPConfig
	Renderers     RendererConfigs
	Kafka         platformkafka.ConnectionConfig
	KafkaConsumer platformkafka.ConsumerConfig
}

type RendererConfig struct {
	TemplatePath string
	ContentType  emailcontract.EmailContentType
}

type RendererConfigs struct {
	Welcome       RendererConfig
	Validation    RendererConfig
	SecurityAlert RendererConfig
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
		Renderers: RendererConfigs{
			Welcome: RendererConfig{
				TemplatePath: "internal/email/templates/welcome_email_template.html",
				ContentType:  emailcontract.EmailContentType_HTML,
			},
			Validation: RendererConfig{
				TemplatePath: "internal/email/templates/validation_email_template.html",
				ContentType:  emailcontract.EmailContentType_HTML,
			},
			SecurityAlert: RendererConfig{
				TemplatePath: "internal/email/templates/security_alert_email_template.html",
				ContentType:  emailcontract.EmailContentType_HTML,
			},
		},
	}
	if config.SMTP.Host == "" || config.SMTP.UserName == "" || config.SMTP.Password == "" || name == "" {
		return Config{}, fmt.Errorf("SMTP_HOST, NOTEZY_OFFICIAL_NAME, NOTEZY_OFFICIAL_GMAIL, and NOTEZY_OFFICIAL_GOOGLE_APPLICATION_PASSWORD are required")
	}
	kafka, err := platformkafka.LoadConnectionConfig()
	if err != nil {
		return Config{}, err
	}
	maximumAttempts, err := positiveIntEnv("KAFKA_CONSUMER_MAXIMUM_ATTEMPTS", 3)
	if err != nil {
		return Config{}, err
	}
	initialRetryBackoff, err := positiveDurationEnv("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF", time.Second)
	if err != nil {
		return Config{}, err
	}
	maximumRetryBackoff, err := positiveDurationEnv("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF", 30*time.Second)
	if err != nil {
		return Config{}, err
	}
	maximumPollRecords, err := positiveIntEnv("KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS", 100)
	if err != nil {
		return Config{}, err
	}
	config.Kafka = kafka
	config.KafkaConsumer = platformkafka.ConsumerConfig{
		ClientConfig: platformkafka.ClientConfig{
			ConnectionConfig: kafka,
			ClientId:         "notezy-email",
		},
		ConsumerGroup:       "notezy-email-core-v1",
		MaximumAttempts:     maximumAttempts,
		InitialRetryBackoff: initialRetryBackoff,
		MaximumRetryBackoff: maximumRetryBackoff,
		MaximumPollRecords:  maximumPollRecords,
	}

	return config, nil
}

func positiveIntEnv(name string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return parsed, nil
}

func positiveDurationEnv(name string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive Go duration", name)
	}
	return parsed, nil
}
