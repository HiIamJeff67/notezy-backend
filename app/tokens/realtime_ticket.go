package tokens

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

func GenerateRealtimeConnectionTicket(
	userPublicId uuid.UUID,
	userAgent string,
) (*string, time.Time, *exceptions.Exception) {
	privateKey, exception := getRealtimeTicketPrivateKey()
	if exception != nil {
		return nil, time.Time{}, exception
	}

	now := time.Now()
	expiresAt := now.Add(constants.RealtimeConnectionTicketExpiresIn)
	userAgentHash := sha256.Sum256([]byte(userAgent))
	claims := types.RealtimeConnectionTicketClaims{
		UserAgentHash:           fmt.Sprintf("%x", userAgentHash),
		RealtimeProtocolVersion: constants.RealtimeProtocolVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{types.RealtimeTicketAudience_Connection.String()},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    constants.ServiceName,
			Subject:   userPublicId.String(),
		},
	}

	ticket, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(privateKey)
	if err != nil {
		return nil, time.Time{}, exceptions.Token.FailedToGenerateRealtimeTicket().WithOrigin(err)
	}

	return &ticket, expiresAt, nil
}

func GenerateRealtimeBlockPackTicket(
	userPublicId uuid.UUID,
	userAgent string,
	blockPackId uuid.UUID,
	permission realtimetypes.ChannelPermission,
) (*string, time.Time, *exceptions.Exception) {
	privateKey, exception := getRealtimeTicketPrivateKey()
	if exception != nil {
		return nil, time.Time{}, exception
	}

	now := time.Now()
	expiresAt := now.Add(constants.RealtimeBlockPackTicketExpiresIn)
	userAgentHash := sha256.Sum256([]byte(userAgent))
	claims := types.RealtimeBlockPackTicketClaims{
		UserAgentHash:           fmt.Sprintf("%x", userAgentHash),
		ChannelType:             string(realtimetypes.ChannelType_BlockPack),
		ChannelId:               blockPackId.String(),
		Permission:              string(permission),
		RealtimeProtocolVersion: constants.RealtimeProtocolVersion,
		SchemaVersion:           constants.YjsBlockPackSchemaVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{types.RealtimeTicketAudience_BlockPack.String()},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    constants.ServiceName,
			Subject:   userPublicId.String(),
		},
	}

	ticket, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims).SignedString(privateKey)
	if err != nil {
		return nil, time.Time{}, exceptions.Token.FailedToGenerateRealtimeTicket().WithOrigin(err)
	}

	return &ticket, expiresAt, nil
}

func ParseRealtimeConnectionTicket(
	ticketString string, userAgent string,
) (*types.RealtimeConnectionTicketClaims, error) {
	privateKey, exception := getRealtimeTicketPrivateKey()
	if exception != nil {
		return nil, fmt.Errorf("%s", exception.Message)
	}

	claims := types.RealtimeConnectionTicketClaims{}
	ticket, err := jwt.ParseWithClaims(
		ticketString,
		&claims,
		func(ticket *jwt.Token) (any, error) {
			if ticket.Method != jwt.SigningMethodEdDSA {
				return nil, fmt.Errorf("invalid realtime ticket signing method")
			}

			return privateKey.Public(), nil
		},
		jwt.WithAudience(types.RealtimeTicketAudience_Connection.String()),
		jwt.WithIssuer(constants.ServiceName),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid realtime connection ticket: %w", err)
	}
	if !ticket.Valid {
		return nil, fmt.Errorf("invalid realtime connection ticket")
	}

	userAgentHash := sha256.Sum256([]byte(userAgent))
	if claims.UserAgentHash != fmt.Sprintf("%x", userAgentHash) ||
		claims.RealtimeProtocolVersion != constants.RealtimeProtocolVersion {
		return nil, fmt.Errorf("invalid realtime connection ticket claims")
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return nil, fmt.Errorf("invalid realtime connection ticket user public id: %w", err)
	}

	return &claims, nil
}

func ParseRealtimeBlockPackTicket(
	ticketString string, userAgent string,
) (*types.RealtimeBlockPackTicketClaims, error) {
	privateKey, exception := getRealtimeTicketPrivateKey()
	if exception != nil {
		return nil, fmt.Errorf("%s", exception.Message)
	}

	claims := types.RealtimeBlockPackTicketClaims{}
	ticket, err := jwt.ParseWithClaims(
		ticketString,
		&claims,
		func(ticket *jwt.Token) (any, error) {
			if ticket.Method != jwt.SigningMethodEdDSA {
				return nil, fmt.Errorf("invalid realtime ticket signing method")
			}

			return privateKey.Public(), nil
		},
		jwt.WithAudience(types.RealtimeTicketAudience_BlockPack.String()),
		jwt.WithIssuer(constants.ServiceName),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid realtime block pack ticket: %w", err)
	}
	if !ticket.Valid {
		return nil, fmt.Errorf("invalid realtime block pack ticket")
	}

	userAgentHash := sha256.Sum256([]byte(userAgent))
	if claims.UserAgentHash != fmt.Sprintf("%x", userAgentHash) ||
		claims.ChannelType != string(realtimetypes.ChannelType_BlockPack) ||
		claims.RealtimeProtocolVersion != constants.RealtimeProtocolVersion ||
		claims.SchemaVersion != constants.YjsBlockPackSchemaVersion {
		return nil, fmt.Errorf("invalid realtime block pack ticket claims")
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return nil, fmt.Errorf("invalid realtime block pack ticket user public id: %w", err)
	}
	if _, err := uuid.Parse(claims.ChannelId); err != nil {
		return nil, fmt.Errorf("invalid realtime block pack ticket channel id: %w", err)
	}
	if claims.Permission != string(realtimetypes.ChannelPermission_Read) &&
		claims.Permission != string(realtimetypes.ChannelPermission_Write) {
		return nil, fmt.Errorf("invalid realtime block pack ticket permission")
	}

	return &claims, nil
}

func getRealtimeTicketPrivateKey() (ed25519.PrivateKey, *exceptions.Exception) {
	encodedPrivateKey := util.GetEnv("REALTIME_TICKET_PRIVATE_KEY_BASE64", "")
	if encodedPrivateKey == "" {
		return nil, exceptions.Token.RealtimeTicketPrivateKeyNotFound()
	}

	privateKeyBytes, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		return nil, exceptions.Token.InvalidRealtimeTicketPrivateKey().WithOrigin(err)
	}

	parsedPrivateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, exceptions.Token.InvalidRealtimeTicketPrivateKey().WithOrigin(err)
	}
	privateKey, ok := parsedPrivateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, exceptions.Token.InvalidRealtimeTicketPrivateKey().WithOrigin(
			fmt.Errorf("expected an Ed25519 PKCS#8 private key"),
		)
	}

	return privateKey, nil
}
