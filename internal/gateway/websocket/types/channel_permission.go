package realtimetypes

import sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"

type ChannelPermission = sharedtypes.ChannelPermission

const (
	ChannelPermission_Read  = sharedtypes.ChannelPermission_Read
	ChannelPermission_Write = sharedtypes.ChannelPermission_Write
)
