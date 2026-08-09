package eventscontract

import "time"

type UserSessionsRevokedData struct{}

type UserDeletedData struct {
	DeletedAt time.Time `json:"deletedAt"`
}
