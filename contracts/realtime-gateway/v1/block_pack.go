package realtimegatewaycontract

import "github.com/google/uuid"

const GetBlockPackParticipantsOperation = "realtimegateway.block-pack-participants.get"

type BlockPackParticipantResponseDto struct {
	UserPublicId      uuid.UUID `json:"userPublicId"`
	ChannelPermission string    `json:"channelPermission"`
	ConnectionCount   int       `json:"connectionCount"`
}

type GetBlockPackParticipantsResponseDto struct {
	Participants []BlockPackParticipantResponseDto `json:"participants"`
}
