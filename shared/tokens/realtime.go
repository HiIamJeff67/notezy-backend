package tokens

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const realtimeConnectionTicketAudience string = "notezy-realtime-connection"
const realtimeBlockPackTicketAudience string = "notezy-realtime-block-pack"
const realtimeTicketIssuer string = "github.com/HiIamJeff67/notezy-backend"
const realtimeConnectionTicketExpiresIn time.Duration = 5 * time.Minute
const realtimeBlockPackTicketExpiresIn time.Duration = 5 * time.Minute
const realtimeProtocolVersion int = 1
const yjsBlockPackSchemaVersion int = 1
const blockPackRoomAdmissionPolicyVersion int = 1
const blockPackDocumentQuotaPolicyVersion int = 1
const roomAdmissionEnforcementStrategyRejectNewSubscriber string = "reject-new-subscriber"

type RealtimeConnectionTicketClaims struct {
	UserAgentHash           string `json:"userAgentHash" validate:"required"`
	RealtimeProtocolVersion int    `json:"realtimeProtocolVersion" validate:"required"`
	jwt.RegisteredClaims
}

type RealtimeBlockPackTicketClaims struct {
	UserAgentHash                    string `json:"userAgentHash" validate:"required"`
	ChannelType                      string `json:"channelType" validate:"required"`
	ChannelId                        string `json:"channelId" validate:"required,uuid4"`
	Permission                       string `json:"permission" validate:"required,oneof=read write"`
	RealtimeProtocolVersion          int    `json:"realtimeProtocolVersion" validate:"required"`
	SchemaVersion                    int    `json:"schemaVersion" validate:"required"`
	RoomAdmissionPolicyVersion       int    `json:"roomAdmissionPolicyVersion" validate:"required"`
	RoomAdmissionEnforcementStrategy string `json:"roomAdmissionEnforcementStrategy" validate:"required"`
	MaximumSubscribers               int32  `json:"maximumSubscribers" validate:"required,min=1"`
	DocumentQuotaPolicyVersion       int    `json:"documentQuotaPolicyVersion" validate:"required"`
	MaximumBlockCount                int32  `json:"maximumBlockCount" validate:"required,min=1"`
	jwt.RegisteredClaims
}

func GenerateRealtimeConnectionTicket(claims RealtimeConnectionTicketClaims) (*string, time.Time, error) {
	privateKey, err := parseRealtimeTicketPrivateKey(os.Getenv("REALTIME_TICKET_PRIVATE_KEY_BASE64"))
	if err != nil {
		return nil, time.Time{}, err
	}
	if claims.UserAgentHash == "" || claims.RealtimeProtocolVersion != realtimeProtocolVersion {
		return nil, time.Time{}, errors.New("realtime connection ticket claims are invalid")
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return nil, time.Time{}, fmt.Errorf("invalid realtime connection ticket user public id: %w", err)
	}

	now := time.Now()
	expiresAt := now.Add(realtimeConnectionTicketExpiresIn)
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{realtimeConnectionTicketAudience},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        uuid.NewString(),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    realtimeTicketIssuer,
		Subject:   claims.Subject,
	}

	ticket, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(privateKey)
	if err != nil {
		return nil, time.Time{}, err
	}

	return &ticket, expiresAt, nil
}

func GenerateRealtimeBlockPackTicket(claims RealtimeBlockPackTicketClaims) (*string, time.Time, error) {
	privateKey, err := parseRealtimeTicketPrivateKey(os.Getenv("REALTIME_TICKET_PRIVATE_KEY_BASE64"))
	if err != nil {
		return nil, time.Time{}, err
	}
	if claims.UserAgentHash == "" ||
		claims.ChannelType != "BlockPack" ||
		claims.RealtimeProtocolVersion != realtimeProtocolVersion ||
		claims.SchemaVersion != yjsBlockPackSchemaVersion ||
		claims.RoomAdmissionPolicyVersion != blockPackRoomAdmissionPolicyVersion ||
		claims.RoomAdmissionEnforcementStrategy != roomAdmissionEnforcementStrategyRejectNewSubscriber ||
		claims.MaximumSubscribers <= 0 ||
		claims.DocumentQuotaPolicyVersion != blockPackDocumentQuotaPolicyVersion ||
		claims.MaximumBlockCount <= 0 ||
		(claims.Permission != "read" && claims.Permission != "write") {
		return nil, time.Time{}, errors.New("realtime block pack ticket claims are invalid")
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return nil, time.Time{}, fmt.Errorf("invalid realtime block pack ticket user public id: %w", err)
	}
	if _, err := uuid.Parse(claims.ChannelId); err != nil {
		return nil, time.Time{}, fmt.Errorf("invalid realtime block pack ticket channel id: %w", err)
	}

	now := time.Now()
	expiresAt := now.Add(realtimeBlockPackTicketExpiresIn)
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{realtimeBlockPackTicketAudience},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        uuid.NewString(),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    realtimeTicketIssuer,
		Subject:   claims.Subject,
	}

	ticket, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(privateKey)
	if err != nil {
		return nil, time.Time{}, err
	}

	return &ticket, expiresAt, nil
}

func ParseRealtimeConnectionTicket(ticketString string, userAgent string) (*RealtimeConnectionTicketClaims, error) {
	publicKey, err := parseRealtimeTicketPublicKey(os.Getenv("REALTIME_TICKET_PUBLIC_KEY_BASE64"))
	if err != nil {
		return nil, err
	}

	claims := RealtimeConnectionTicketClaims{}
	ticket, err := jwt.ParseWithClaims(
		ticketString,
		&claims,
		func(ticket *jwt.Token) (any, error) {
			if ticket.Method != jwt.SigningMethodEdDSA {
				return nil, jwt.ErrSignatureInvalid
			}

			return publicKey, nil
		},
		jwt.WithAudience(realtimeConnectionTicketAudience),
		jwt.WithIssuer(realtimeTicketIssuer),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid realtime connection ticket: %w", err)
	}
	if !ticket.Valid {
		return nil, fmt.Errorf("invalid realtime connection ticket")
	}

	userAgentHash := sha256.Sum256([]byte(userAgent))
	if claims.ID == "" ||
		claims.ExpiresAt == nil ||
		claims.UserAgentHash != fmt.Sprintf("%x", userAgentHash) ||
		claims.RealtimeProtocolVersion != realtimeProtocolVersion {
		return nil, fmt.Errorf("invalid realtime connection ticket claims")
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return nil, fmt.Errorf("invalid realtime connection ticket user public id: %w", err)
	}

	return &claims, nil
}

func ParseRealtimeBlockPackTicket(ticketString string, userAgent string) (*RealtimeBlockPackTicketClaims, error) {
	publicKey, err := parseRealtimeTicketPublicKey(os.Getenv("REALTIME_TICKET_PUBLIC_KEY_BASE64"))
	if err != nil {
		return nil, err
	}

	claims := RealtimeBlockPackTicketClaims{}
	ticket, err := jwt.ParseWithClaims(
		ticketString,
		&claims,
		func(ticket *jwt.Token) (any, error) {
			if ticket.Method != jwt.SigningMethodEdDSA {
				return nil, jwt.ErrSignatureInvalid
			}

			return publicKey, nil
		},
		jwt.WithAudience(realtimeBlockPackTicketAudience),
		jwt.WithIssuer(realtimeTicketIssuer),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid realtime block pack ticket: %w", err)
	}
	if !ticket.Valid {
		return nil, fmt.Errorf("invalid realtime block pack ticket")
	}

	userAgentHash := sha256.Sum256([]byte(userAgent))
	if claims.ID == "" ||
		claims.ExpiresAt == nil ||
		claims.UserAgentHash != fmt.Sprintf("%x", userAgentHash) ||
		claims.ChannelType != "BlockPack" ||
		claims.RealtimeProtocolVersion != realtimeProtocolVersion ||
		claims.SchemaVersion != yjsBlockPackSchemaVersion ||
		claims.RoomAdmissionPolicyVersion != blockPackRoomAdmissionPolicyVersion ||
		claims.RoomAdmissionEnforcementStrategy != roomAdmissionEnforcementStrategyRejectNewSubscriber ||
		claims.MaximumSubscribers <= 0 ||
		claims.DocumentQuotaPolicyVersion != blockPackDocumentQuotaPolicyVersion ||
		claims.MaximumBlockCount <= 0 {
		return nil, fmt.Errorf("invalid realtime block pack ticket claims")
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return nil, fmt.Errorf("invalid realtime block pack ticket user public id: %w", err)
	}
	if _, err := uuid.Parse(claims.ChannelId); err != nil {
		return nil, fmt.Errorf("invalid realtime block pack ticket channel id: %w", err)
	}
	if claims.Permission != "read" && claims.Permission != "write" {
		return nil, fmt.Errorf("invalid realtime block pack ticket permission")
	}

	return &claims, nil
}

func parseRealtimeTicketPrivateKey(encodedPrivateKey string) (ed25519.PrivateKey, error) {
	if encodedPrivateKey == "" {
		return nil, errors.New("realtime ticket private key is required")
	}

	privateKeyBytes, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid realtime ticket private key: %w", err)
	}
	parsedPrivateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid realtime ticket private key: %w", err)
	}
	privateKey, ok := parsedPrivateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("invalid realtime ticket private key")
	}

	return privateKey, nil
}

func parseRealtimeTicketPublicKey(encodedPublicKey string) (ed25519.PublicKey, error) {
	if encodedPublicKey == "" {
		return nil, errors.New("realtime ticket public key is required")
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(encodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid realtime ticket public key: %w", err)
	}
	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid realtime ticket public key: %w", err)
	}
	publicKey, ok := parsedPublicKey.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("invalid realtime ticket public key")
	}

	return publicKey, nil
}
