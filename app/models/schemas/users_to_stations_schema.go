package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type UsersToStations struct {
	UserId     uuid.UUID                     `json:"userId" gorm:"column:user_id; type:uuid; primaryKey;"`
	StationId  uuid.UUID                     `json:"stationId" gorm:"column:station_id; type:uuid; primaryKey;"`
	Permission enums.AccessControlPermission `json:"permission" gorm:"column:permission; type:\"AccessControlPermission\"; not null; default:'Read';"`
	UpdatedAt  time.Time                     `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt  time.Time                     `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	User    User    `json:"user" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Station Station `json:"station" gorm:"foreignKey:StationId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// UsersToStations Table Name
func (UsersToStations) TableName() string {
	return types.TableName_UsersToStationsTable.String()
}

// UsersToStations Table Relations
type UsersToStationsRelation types.RelationName

const (
	UsersToStationsRelation_User    UsersToStationsRelation = "User"
	UsersToStationsRelation_Station UsersToStationsRelation = "Station"
)
