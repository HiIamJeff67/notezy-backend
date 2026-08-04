package routinetask

import (
	"testing"

	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	durablejobconfig "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/config"
)

func TestNewHandlerManagerRegistersEveryPurposePolicy(t *testing.T) {
	manager := NewHandlerManager(1, nil, durablejobconfig.Config{})

	for _, purpose := range enums.AllRoutineTaskPurposes {
		registry, exists := manager.registries[purpose]
		if !exists || registry.HandlerFunc == nil || len(registry.AllowedPermissions) == 0 {
			t.Fatalf("missing routine task handler policy for %s", purpose)
		}
	}
}
