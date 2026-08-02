package configs

import "os"

func RealtimeGatewayListenAddress() string {
	domain := os.Getenv("REALTIME_GATEWAY_DOMAIN")
	if domain == "" {
		domain = os.Getenv("REALTIME_DOMAIN")
	}
	if domain == "" {
		domain = os.Getenv("WEBSOCKET_DOMAIN")
	}
	if domain == "" {
		domain = "0.0.0.0"
	}

	port := os.Getenv("REALTIME_GATEWAY_PORT")
	if port == "" {
		port = os.Getenv("REALTIME_PORT")
	}
	if port == "" {
		port = os.Getenv("WEBSOCKET_PORT")
	}
	if port == "" {
		port = "7779"
	}

	return domain + ":" + port
}
