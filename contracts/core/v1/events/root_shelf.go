package eventscontract

import "github.com/google/uuid"

type RootShelfPermissionRevokedData struct {
	TargetUserPublicId *uuid.UUID `json:"targetUserPublicId,omitempty"`
}
