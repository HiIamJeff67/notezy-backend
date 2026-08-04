package eventscontract

import "github.com/google/uuid"

type BlockPackAccessRevocationReason string

const (
	BlockPackAccessRevocationReason_PermissionRevoked   BlockPackAccessRevocationReason = "permission_revoked"
	BlockPackAccessRevocationReason_ResourceUnavailable BlockPackAccessRevocationReason = "resource_unavailable"
)

type BlockPackAccessRevokedData struct {
	TargetUserPublicId *uuid.UUID                      `json:"targetUserPublicId,omitempty"`
	Reason             BlockPackAccessRevocationReason `json:"reason"`
}

type BlockPackRoomPolicyChangedData struct {
	RoomAdmissionPolicyVersion       int                              `json:"roomAdmissionPolicyVersion"`
	RoomAdmissionEnforcementStrategy RoomAdmissionEnforcementStrategy `json:"roomAdmissionEnforcementStrategy"`
	MaximumSubscribers               int32                            `json:"maximumSubscribers"`
}
