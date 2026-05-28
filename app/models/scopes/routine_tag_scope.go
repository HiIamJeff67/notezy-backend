package scopes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
)

type RoutineTagScopeInterface interface {
	PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	IncludePreloads(preloads []schemas.RoutineTagRelation) func(db *gorm.DB) *gorm.DB
}

type RoutineTagScope struct{}

func NewRoutineTagScope() RoutineTagScopeInterface {
	return &RoutineTagScope{}
}

func (sc *RoutineTagScope) PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Current schema keeps RoutineTag owner-scoped, while shared visibility comes through assigned routines.
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToStations{}).
			Select("1").
			Joins("INNER JOIN \"RoutineTable\" r ON r.station_id = \"UsersToStationsTable\".station_id").
			Joins("INNER JOIN \"RoutinesToTagsTable\" rtt ON rtt.routine_id = r.id").
			Where("rtt.tag_id = \"RoutineTagTable\".id").
			Where("\"UsersToStationsTable\".user_id = ? AND \"UsersToStationsTable\".permission IN ?", userId, permissions)
		return db.Where("\"RoutineTagTable\".id = ? AND (\"RoutineTagTable\".owner_id = ? OR EXISTS (?))", id, userId, subQuery)
	}
}

func (sc *RoutineTagScope) PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Current schema keeps RoutineTag owner-scoped, while shared visibility comes through assigned routines.
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToStations{}).
			Select("1").
			Joins("INNER JOIN \"RoutineTable\" r ON r.station_id = \"UsersToStationsTable\".station_id").
			Joins("INNER JOIN \"RoutinesToTagsTable\" rtt ON rtt.routine_id = r.id").
			Where("rtt.tag_id = \"RoutineTagTable\".id").
			Where("\"UsersToStationsTable\".user_id = ? AND \"UsersToStationsTable\".permission IN ?", userId, permissions)
		return db.Where("\"RoutineTagTable\".id IN ? AND (\"RoutineTagTable\".owner_id = ? OR EXISTS (?))", ids, userId, subQuery)
	}
}

func (sc *RoutineTagScope) IncludePreloads(preloads []schemas.RoutineTagRelation) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
		return db
	}
}
