package tokens

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

func TestGenerateRealtimeConnectionTicket(t *testing.T) {
	publicKey := configureRealtimeTicketPrivateKey(t)
	userPublicId := uuid.New()

	ticket, expiresAt, exception := GenerateRealtimeConnectionTicket(userPublicId, "test-user-agent")
	if exception != nil {
		t.Fatalf("failed to generate connection ticket: %v", exception)
	}
	if ticket == nil || expiresAt.IsZero() {
		t.Fatalf("expected a ticket and expiry")
	}

	claims := types.RealtimeConnectionTicketClaims{}
	parsedTicket, err := jwt.ParseWithClaims(*ticket, &claims, func(ticket *jwt.Token) (any, error) {
		if ticket.Method != jwt.SigningMethodEdDSA {
			return nil, jwt.ErrSignatureInvalid
		}

		return publicKey, nil
	}, jwt.WithAudience(types.RealtimeTicketAudience_Connection.String()))
	if err != nil || !parsedTicket.Valid {
		t.Fatalf("failed to verify connection ticket: %v", err)
	}
	if claims.Subject != userPublicId.String() || claims.RealtimeProtocolVersion != constants.RealtimeProtocolVersion {
		t.Fatalf("unexpected connection claims: %#v", claims)
	}

}

func TestGenerateRealtimeBlockPackTicket(t *testing.T) {
	publicKey := configureRealtimeTicketPrivateKey(t)
	userPublicId := uuid.New()
	blockPackId := uuid.New()

	ticket, expiresAt, exception := GenerateRealtimeBlockPackTicket(
		userPublicId,
		"test-user-agent",
		blockPackId,
		realtimetypes.ChannelPermission_Write,
	)
	if exception != nil {
		t.Fatalf("failed to generate block pack ticket: %v", exception)
	}
	if ticket == nil || expiresAt.IsZero() {
		t.Fatalf("expected a ticket and expiry")
	}

	claims := types.RealtimeBlockPackTicketClaims{}
	parsedTicket, err := jwt.ParseWithClaims(*ticket, &claims, func(ticket *jwt.Token) (any, error) {
		if ticket.Method != jwt.SigningMethodEdDSA {
			return nil, jwt.ErrSignatureInvalid
		}

		return publicKey, nil
	}, jwt.WithAudience(types.RealtimeTicketAudience_BlockPack.String()))
	if err != nil || !parsedTicket.Valid {
		t.Fatalf("failed to verify block pack ticket: %v", err)
	}
	if claims.Subject != userPublicId.String() ||
		claims.ChannelId != blockPackId.String() ||
		claims.ChannelType != string(realtimetypes.ChannelType_BlockPack) ||
		claims.Permission != string(realtimetypes.ChannelPermission_Write) ||
		claims.SchemaVersion != constants.YjsBlockPackSchemaVersion {
		t.Fatalf("unexpected block pack claims: %#v", claims)
	}
}

func configureRealtimeTicketPrivateKey(t *testing.T) ed25519.PublicKey {
	t.Helper()

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to marshal test private key: %v", err)
	}

	t.Setenv("REALTIME_TICKET_PRIVATE_KEY_BASE64", base64.StdEncoding.EncodeToString(privateKeyBytes))

	return publicKey
}
