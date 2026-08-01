package types

import "fmt"

type AccessControlPermission string

const (
	AccessControlPermission_Read  AccessControlPermission = "Read"
	AccessControlPermission_Write AccessControlPermission = "Write"
	AccessControlPermission_Admin AccessControlPermission = "Admin"
	AccessControlPermission_Owner AccessControlPermission = "Owner"
)

var AllAccessControlPermissions = []AccessControlPermission{
	AccessControlPermission_Read,
	AccessControlPermission_Write,
	AccessControlPermission_Admin,
	AccessControlPermission_Owner,
}

func (p AccessControlPermission) String() string {
	return string(p)
}

func ParseAccessControlPermission(value string) (*AccessControlPermission, error) {
	for _, permission := range AllAccessControlPermissions {
		if permission == AccessControlPermission(value) {
			return &permission, nil
		}
	}

	return nil, fmt.Errorf("invalid access control permission: %s", value)
}
