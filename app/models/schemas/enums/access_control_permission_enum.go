package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

/* ============================== Definition ============================== */

type AccessControlPermission string

const (
	AccessControlPermission_Read  AccessControlPermission = "Read"
	AccessControlPermission_Write AccessControlPermission = "Write"
	AccessControlPermission_Admin AccessControlPermission = "Admin"
)

/* ============================== All Instances ============================== */

var AllAccessControlPermissions = []AccessControlPermission{
	AccessControlPermission_Read,
	AccessControlPermission_Write,
	AccessControlPermission_Admin,
}

var AllAccessControlPermissionStrings = []string{
	string(AccessControlPermission_Read),
	string(AccessControlPermission_Write),
	string(AccessControlPermission_Admin),
}

/* ============================== Methods ============================== */

func (a AccessControlPermission) Name() string {
	return reflect.TypeOf(a).Name()
}

func (a *AccessControlPermission) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*a = AccessControlPermission(string(v))
		return nil
	case string:
		*a = AccessControlPermission(v)
		return nil
	}
	return scanError(value, a)
}

func (a AccessControlPermission) Value() (driver.Value, error) {
	return string(a), nil
}

func (a AccessControlPermission) String() string {
	return string(a)
}

func (a *AccessControlPermission) IsValidEnum() bool {
	return slices.Contains(AllAccessControlPermissions, *a)
}

func ConvertStringToAccessControlPermission(enumString string) (*AccessControlPermission, error) {
	for _, accessControlPermission := range AllAccessControlPermissions {
		if string(accessControlPermission) == enumString {
			return &accessControlPermission, nil
		}
	}
	return nil, fmt.Errorf("invalid access control permission: %s", enumString)
}
