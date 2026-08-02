package configs

import "os"

func GatewayListenAddress() string {
	domain := os.Getenv("GIN_DOMAIN")
	if domain == "" {
		domain = "0.0.0.0"
	}
	port := os.Getenv("GIN_PORT")
	if port == "" {
		port = "7777"
	}

	return domain + ":" + port
}
