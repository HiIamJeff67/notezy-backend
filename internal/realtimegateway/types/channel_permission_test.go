package realtimetypes

import (
	"testing"

	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func TestChannelPermissionAllowedAccessControlPermissions(t *testing.T) {
	readPermissions := ChannelPermission_Read.AllowedAccessControlPermissions()
	if len(readPermissions) != 4 || readPermissions[3] != sharedtypes.AccessControlPermission_Read {
		t.Fatalf("expected read channel to include Read access: %#v", readPermissions)
	}

	writePermissions := ChannelPermission_Write.AllowedAccessControlPermissions()
	if len(writePermissions) != 3 || writePermissions[2] != sharedtypes.AccessControlPermission_Write {
		t.Fatalf("expected write channel to stop at Write access: %#v", writePermissions)
	}

	if permissions := ChannelPermission("invalid").AllowedAccessControlPermissions(); permissions != nil {
		t.Fatalf("expected invalid channel permission to have no access control permissions: %#v", permissions)
	}
}
