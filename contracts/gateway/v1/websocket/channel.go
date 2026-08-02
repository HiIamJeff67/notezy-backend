package websocketcontract

import "github.com/google/uuid"

type ValidateBlockPackChannelPermissionRequestDto struct {
	UserPublicId uuid.UUID `json:"userPublicId" validate:"required"`
	BlockPackId  uuid.UUID `json:"blockPackId" validate:"required"`
	Permission   string    `json:"permission" validate:"required,oneof=read write"`
}

type ValidateBlockPackChannelPermissionResponseDto struct {
	Permitted bool   `json:"permitted"`
	ErrorCode string `json:"errorCode"`
}
