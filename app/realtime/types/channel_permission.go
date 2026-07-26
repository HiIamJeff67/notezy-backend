package realtimetypes

import (
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

type ChannelPermission string

const (
	ChannelPermission_Read  ChannelPermission = "read"
	ChannelPermission_Write ChannelPermission = "write"
)

func (p ChannelPermission) AllowedAccessControlPermissions() []enums.AccessControlPermission {
	switch p {
	case ChannelPermission_Read:
		return []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		}
	case ChannelPermission_Write:
		return []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
	default:
		return nil
	}
}
