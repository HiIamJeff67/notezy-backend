package tokens

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestRealtimeConnectionTicketRoundTrip(t *testing.T) {
	t.Setenv("REALTIME_TICKET_PRIVATE_KEY_BASE64", generateRealtimePrivateKey(t))
	userPublicId := uuid.New()
	userAgentHash := sha256.Sum256([]byte("test-agent"))
	claims := RealtimeConnectionTicketClaims{
		UserAgentHash:           fmt.Sprintf("%x", userAgentHash),
		RealtimeProtocolVersion: realtimeProtocolVersion,
	}
	claims.Subject = userPublicId.String()

	ticket, _, err := GenerateRealtimeConnectionTicket(claims)
	if err != nil {
		t.Fatalf("generate realtime connection ticket: %v", err)
	}

	parsedClaims, err := ParseRealtimeConnectionTicket(*ticket, "test-agent")
	if err != nil {
		t.Fatalf("parse realtime connection ticket: %v", err)
	}
	if parsedClaims.Subject != userPublicId.String() {
		t.Fatalf("unexpected connection ticket claims: %#v", parsedClaims)
	}
}

func TestRealtimeBlockPackTicketRoundTrip(t *testing.T) {
	t.Setenv("REALTIME_TICKET_PRIVATE_KEY_BASE64", generateRealtimePrivateKey(t))
	userPublicId := uuid.New()
	blockPackId := uuid.New()
	userAgentHash := sha256.Sum256([]byte("test-agent"))
	claims := RealtimeBlockPackTicketClaims{
		UserAgentHash:                    fmt.Sprintf("%x", userAgentHash),
		ChannelType:                      "BlockPack",
		ChannelId:                        blockPackId.String(),
		Permission:                       "write",
		RealtimeProtocolVersion:          realtimeProtocolVersion,
		SchemaVersion:                    yjsBlockPackSchemaVersion,
		RoomAdmissionPolicyVersion:       blockPackRoomAdmissionPolicyVersion,
		RoomAdmissionEnforcementStrategy: roomAdmissionEnforcementStrategyRejectNewSubscriber,
		MaximumSubscribers:               10,
	}
	claims.Subject = userPublicId.String()

	ticket, _, err := GenerateRealtimeBlockPackTicket(claims)
	if err != nil {
		t.Fatalf("generate realtime block pack ticket: %v", err)
	}

	parsedClaims, err := ParseRealtimeBlockPackTicket(*ticket, "test-agent")
	if err != nil {
		t.Fatalf("parse realtime block pack ticket: %v", err)
	}
	if parsedClaims.Subject != userPublicId.String() ||
		parsedClaims.ChannelId != blockPackId.String() ||
		parsedClaims.RoomAdmissionPolicyVersion != blockPackRoomAdmissionPolicyVersion ||
		parsedClaims.RoomAdmissionEnforcementStrategy != roomAdmissionEnforcementStrategyRejectNewSubscriber ||
		parsedClaims.MaximumSubscribers != 10 {
		t.Fatalf("unexpected block pack ticket claims: %#v", parsedClaims)
	}
}

func TestGenerateRealtimeBlockPackTicketRejectsUnsupportedRoomAdmissionPolicy(t *testing.T) {
	t.Setenv("REALTIME_TICKET_PRIVATE_KEY_BASE64", generateRealtimePrivateKey(t))
	userPublicId := uuid.New()
	blockPackId := uuid.New()
	userAgentHash := sha256.Sum256([]byte("test-agent"))
	claims := RealtimeBlockPackTicketClaims{
		UserAgentHash:                    fmt.Sprintf("%x", userAgentHash),
		ChannelType:                      "BlockPack",
		ChannelId:                        blockPackId.String(),
		Permission:                       "write",
		RealtimeProtocolVersion:          realtimeProtocolVersion,
		SchemaVersion:                    yjsBlockPackSchemaVersion,
		RoomAdmissionPolicyVersion:       blockPackRoomAdmissionPolicyVersion + 1,
		RoomAdmissionEnforcementStrategy: roomAdmissionEnforcementStrategyRejectNewSubscriber,
		MaximumSubscribers:               10,
	}
	claims.Subject = userPublicId.String()

	if _, _, err := GenerateRealtimeBlockPackTicket(claims); err == nil {
		t.Fatal("expected unsupported room admission policy to be rejected")
	}
}

func generateRealtimePrivateKey(t *testing.T) string {
	t.Helper()

	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate private key: %v", err)
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal private key: %v", err)
	}

	return base64.StdEncoding.EncodeToString(privateKeyBytes)
}
