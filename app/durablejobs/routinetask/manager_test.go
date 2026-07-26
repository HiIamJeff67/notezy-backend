package routinetask

import (
	"testing"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

func TestNewHandlerManagerRegistersEveryPurposePolicy(t *testing.T) {
	manager := NewHandlerManager(1, nil)

	for _, purpose := range enums.AllRoutineTaskPurposes {
		registry, exists := manager.registries[purpose]
		if !exists || registry.HandlerFunc == nil || len(registry.AllowedPermissions) == 0 {
			t.Fatalf("missing routine task handler policy for %s", purpose)
		}
	}
}
