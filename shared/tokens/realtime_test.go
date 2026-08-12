package tokens

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestRealtimeConnectionTicketRoundTrip(t *testing.T) {
	configureRealtimeTicketKeys(t)
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

func TestRealtimeConnectionTicketRejectsAnotherPublicKey(t *testing.T) {
	configureRealtimeTicketKeys(t)
	userAgentHash := sha256.Sum256([]byte("test-agent"))
	claims := RealtimeConnectionTicketClaims{
		UserAgentHash:           fmt.Sprintf("%x", userAgentHash),
		RealtimeProtocolVersion: realtimeProtocolVersion,
	}
	claims.Subject = uuid.NewString()
	ticket, _, err := GenerateRealtimeConnectionTicket(claims)
	if err != nil {
		t.Fatalf("generate realtime connection ticket: %v", err)
	}

	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate another public key: %v", err)
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("marshal another public key: %v", err)
	}
	t.Setenv("REALTIME_TICKET_PUBLIC_KEY_BASE64", base64.StdEncoding.EncodeToString(publicKeyBytes))

	if _, err := ParseRealtimeConnectionTicket(*ticket, "test-agent"); err == nil {
		t.Fatal("expected another public key to reject the realtime connection ticket")
	}
}

func TestRealtimeConnectionTicketRejectsModifiedClaims(t *testing.T) {
	configureRealtimeTicketKeys(t)
	userAgentHash := sha256.Sum256([]byte("test-agent"))
	claims := RealtimeConnectionTicketClaims{
		UserAgentHash:           fmt.Sprintf("%x", userAgentHash),
		RealtimeProtocolVersion: realtimeProtocolVersion,
	}
	claims.Subject = uuid.NewString()
	ticket, _, err := GenerateRealtimeConnectionTicket(claims)
	if err != nil {
		t.Fatalf("generate realtime connection ticket: %v", err)
	}

	segments := strings.Split(*ticket, ".")
	if len(segments) != 3 {
		t.Fatal("unexpected realtime ticket format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(segments[1])
	if err != nil {
		t.Fatalf("decode realtime ticket claims: %v", err)
	}
	var modifiedClaims map[string]any
	if err := json.Unmarshal(payload, &modifiedClaims); err != nil {
		t.Fatalf("unmarshal realtime ticket claims: %v", err)
	}
	modifiedClaims["sub"] = uuid.NewString()
	payload, err = json.Marshal(modifiedClaims)
	if err != nil {
		t.Fatalf("marshal modified realtime ticket claims: %v", err)
	}
	segments[1] = base64.RawURLEncoding.EncodeToString(payload)
	modifiedTicket := strings.Join(segments, ".")

	if _, err := ParseRealtimeConnectionTicket(modifiedTicket, "test-agent"); err == nil {
		t.Fatal("expected modified realtime ticket claims to be rejected")
	}
}

func TestRealtimeTicketRequiresRuntimeSpecificKey(t *testing.T) {
	configureRealtimeTicketKeys(t)
	userAgentHash := sha256.Sum256([]byte("test-agent"))
	claims := RealtimeConnectionTicketClaims{
		UserAgentHash:           fmt.Sprintf("%x", userAgentHash),
		RealtimeProtocolVersion: realtimeProtocolVersion,
	}
	claims.Subject = uuid.NewString()

	t.Setenv("REALTIME_TICKET_PRIVATE_KEY_BASE64", "")
	if _, _, err := GenerateRealtimeConnectionTicket(claims); err == nil {
		t.Fatal("expected a missing realtime ticket private key to reject generation")
	}
	configureRealtimeTicketKeys(t)
	ticket, _, err := GenerateRealtimeConnectionTicket(claims)
	if err != nil {
		t.Fatalf("generate realtime connection ticket: %v", err)
	}
	t.Setenv("REALTIME_TICKET_PUBLIC_KEY_BASE64", "")
	if _, err := ParseRealtimeConnectionTicket(*ticket, "test-agent"); err == nil {
		t.Fatal("expected a missing realtime ticket public key to reject parsing")
	}
}

func TestRealtimeBlockPackTicketRoundTrip(t *testing.T) {
	configureRealtimeTicketKeys(t)
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
		DocumentQuotaPolicyVersion:       blockPackDocumentQuotaPolicyVersion,
		MaximumBlockCount:                1000,
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
		parsedClaims.MaximumSubscribers != 10 ||
		parsedClaims.DocumentQuotaPolicyVersion != blockPackDocumentQuotaPolicyVersion ||
		parsedClaims.MaximumBlockCount != 1000 {
		t.Fatalf("unexpected block pack ticket claims: %#v", parsedClaims)
	}
}

func TestGenerateRealtimeBlockPackTicketRejectsUnsupportedRoomAdmissionPolicy(t *testing.T) {
	configureRealtimeTicketKeys(t)
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
		DocumentQuotaPolicyVersion:       blockPackDocumentQuotaPolicyVersion,
		MaximumBlockCount:                1000,
	}
	claims.Subject = userPublicId.String()

	if _, _, err := GenerateRealtimeBlockPackTicket(claims); err == nil {
		t.Fatal("expected unsupported room admission policy to be rejected")
	}
}

func configureRealtimeTicketKeys(t *testing.T) {
	t.Helper()

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate private key: %v", err)
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal private key: %v", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}

	t.Setenv("REALTIME_TICKET_PRIVATE_KEY_BASE64", base64.StdEncoding.EncodeToString(privateKeyBytes))
	t.Setenv("REALTIME_TICKET_PUBLIC_KEY_BASE64", base64.StdEncoding.EncodeToString(publicKeyBytes))
}
