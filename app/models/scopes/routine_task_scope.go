package scopes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
)

type RoutineTaskScopeInterface interface {
	PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	IncludePreloads(preloads []schemas.RoutineTaskRelation) func(db *gorm.DB) *gorm.DB
}

type RoutineTaskScope struct{}

func NewRoutineTaskScope() RoutineTaskScopeInterface {
	return &RoutineTaskScope{}
}

func (sc *RoutineTaskScope) PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Current schema keeps RoutineTask owner-scoped, while shared visibility comes through assigned routines.
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToStations{}).
			Select("1").
			Joins("INNER JOIN \"RoutineTable\" r ON r.station_id = \"UsersToStationsTable\".station_id").
			Joins("INNER JOIN \"RoutinesToTasksTable\" rtt ON rtt.routine_id = r.id").
			Where("rtt.task_id = \"RoutineTaskTable\".id").
			Where("\"UsersToStationsTable\".user_id = ? AND \"UsersToStationsTable\".permission IN ?", userId, permissions)
		return db.Where("\"RoutineTaskTable\".id = ? AND (\"RoutineTaskTable\".owner_id = ? OR EXISTS (?))", id, userId, subQuery)
	}
}

func (sc *RoutineTaskScope) PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Current schema keeps RoutineTask owner-scoped, while shared visibility comes through assigned routines.
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToStations{}).
			Select("1").
			Joins("INNER JOIN \"RoutineTable\" r ON r.station_id = \"UsersToStationsTable\".station_id").
			Joins("INNER JOIN \"RoutinesToTasksTable\" rtt ON rtt.routine_id = r.id").
			Where("rtt.task_id = \"RoutineTaskTable\".id").
			Where("\"UsersToStationsTable\".user_id = ? AND \"UsersToStationsTable\".permission IN ?", userId, permissions)
		return db.Where("\"RoutineTaskTable\".id IN ? AND (\"RoutineTaskTable\".owner_id = ? OR EXISTS (?))", ids, userId, subQuery)
	}
}

func (sc *RoutineTaskScope) IncludePreloads(preloads []schemas.RoutineTaskRelation) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
		return db
	}
}
