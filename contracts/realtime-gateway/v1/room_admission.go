package realtimegatewaycontract

type RoomAdmissionEnforcementStrategy string

const (
	BlockPackRoomAdmissionPolicyVersion = 1

	RoomAdmissionEnforcementStrategy_RejectNewSubscriber RoomAdmissionEnforcementStrategy = "reject-new-subscriber"
)
