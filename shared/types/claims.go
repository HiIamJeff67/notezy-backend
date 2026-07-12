package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Name      string `json:"name" validate:"required,min=6,max=16,alphaandnum"`
	Email     string `json:"email" validate:"required,email"`
	UserAgent string `json:"userAgent" validate:"required"`
	jwt.RegisteredClaims
}

type CSRFClaims struct {
	Signature string    `json:"signature" validate:"required"`
	ExpiresAt time.Time `json:"expiresAt"`
	IssuedAt  time.Time `json:"issuedAt"`
}

type RealtimeConnectionTicketClaims struct {
	UserAgentHash           string `json:"userAgentHash" validate:"required"`
	RealtimeProtocolVersion int    `json:"realtimeProtocolVersion" validate:"required"`
	jwt.RegisteredClaims
}

type RealtimeBlockPackTicketClaims struct {
	UserAgentHash           string `json:"userAgentHash" validate:"required"`
	ChannelType             string `json:"channelType" validate:"required"`
	ChannelId               string `json:"channelId" validate:"required,uuid4"`
	Permission              string `json:"permission" validate:"required,oneof=read write"`
	RealtimeProtocolVersion int    `json:"realtimeProtocolVersion" validate:"required"`
	SchemaVersion           int    `json:"schemaVersion" validate:"required"`
	jwt.RegisteredClaims
}
