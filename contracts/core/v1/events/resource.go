package eventscontract

import "github.com/google/uuid"

type ResourceEventChange string

const (
	ResourceEventChange_PermissionUpdated ResourceEventChange = "permission_updated"
	ResourceEventChange_PermissionRevoked ResourceEventChange = "permission_revoked"
	ResourceEventChange_Updated           ResourceEventChange = "updated"
	ResourceEventChange_Deleted           ResourceEventChange = "deleted"
)

type ResourceChangedData struct {
	ResourceId         uuid.UUID           `json:"resourceId"`
	TargetUserPublicId *uuid.UUID          `json:"targetUserPublicId,omitempty"`
	Change             ResourceEventChange `json:"change"`
	Permission         string              `json:"permission,omitempty"`
}
