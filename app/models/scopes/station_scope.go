package scopes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type StationScopeInterface interface {
	PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permission []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permission []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	FilterOnlyDeleted(onlyDeleted types.Ternary) func(db *gorm.DB) *gorm.DB
	IncludePreloads(preloads []schemas.StationRelation) func(db *gorm.DB) *gorm.DB
}

type StationScope struct{}

func NewStationScope() StationScopeInterface {
	return &StationScope{}
}

func (ss *StationScope) PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permission []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToStations{}).
			Select("1").
			Where("station_id = \"StationTable\".id AND user_id = ? AND permission IN ?", userId, permission)
		return db.Where("id = ? AND EXISTS (?)", id, subQuery)
	}
}

func (ss *StationScope) PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permission []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToStations{}).
			Select("1").
			Where("station_id = \"StationTable\".id AND user_id = ? AND permission IN ?", userId, permission)
		return db.Where("\"StationTable\".id IN ? AND EXISTS (?)", ids, subQuery)
	}
}

func (ss *StationScope) FilterOnlyDeleted(onlyDeleted types.Ternary) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch onlyDeleted {
		case types.Ternary_Positive:
			return db.Where("\"StationTable\".deleted_at IS NOT NULL")
		case types.Ternary_Negative:
			return db.Where("\"StationTable\".deleted_at IS NULL")
		default:
			return db
		}
	}
}

func (ss *StationScope) IncludePreloads(preloads []schemas.StationRelation) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
		return db
	}
}
