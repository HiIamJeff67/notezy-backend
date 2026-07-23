package contexts

import (
	"context"
	"testing"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

func TestIntersectAllowedPermissions(t *testing.T) {
	ctx := WithAllowedPermissions(
		context.Background(),
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Owner,
		},
	)

	permissions := IntersectAllowedPermissions(
		ctx,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		},
	)

	expectedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}
	if len(permissions) != len(expectedPermissions) {
		t.Fatalf("expected %d permissions, got %d", len(expectedPermissions), len(permissions))
	}
	for index, expectedPermission := range expectedPermissions {
		if permissions[index] != expectedPermission {
			t.Fatalf("expected permission %s at index %d, got %s", expectedPermission, index, permissions[index])
		}
	}
}

func TestIntersectAllowedPermissionsFailsClosedWithoutPolicy(t *testing.T) {
	permissions := IntersectAllowedPermissions(
		context.Background(),
		[]enums.AccessControlPermission{enums.AccessControlPermission_Read},
	)

	if len(permissions) != 0 {
		t.Fatalf("expected no permissions without a context policy, got %d", len(permissions))
	}
}
