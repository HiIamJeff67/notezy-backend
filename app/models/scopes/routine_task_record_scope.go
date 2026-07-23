package scopes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

type RoutineTaskRecordScopeInterface interface {
	PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	IncludePreloads(preloads []schemas.RoutineTaskRecordRelation) func(db *gorm.DB) *gorm.DB
}

type RoutineTaskRecordScope struct{}

func NewRoutineTaskRecordScope() RoutineTaskRecordScopeInterface {
	return &RoutineTaskRecordScope{}
}

func (sc *RoutineTaskRecordScope) PassPermissionCheck(
	id uuid.UUID,
	userId uuid.UUID,
	permissions []enums.AccessControlPermission,
) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		allowedPermissions := contexts.IntersectAllowedPermissions(db.Statement.Context, permissions)

		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.RoutineTask{}).
			Select("1").
			Joins(`INNER JOIN "RoutineTable" routine ON routine.id = "RoutineTaskTable".routine_id AND routine.deleted_at IS NULL`).
			Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = routine.station_id`).
			Where(`"RoutineTaskTable".id = "RoutineTaskRecordTable".routine_task_id`).
			Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions)
		return db.Where(`"RoutineTaskRecordTable".id = ? AND EXISTS (?)`, id, subQuery)
	}
}

func (sc *RoutineTaskRecordScope) PassPermissionChecks(
	ids []uuid.UUID,
	userId uuid.UUID,
	permissions []enums.AccessControlPermission,
) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		allowedPermissions := contexts.IntersectAllowedPermissions(db.Statement.Context, permissions)

		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.RoutineTask{}).
			Select("1").
			Joins(`INNER JOIN "RoutineTable" routine ON routine.id = "RoutineTaskTable".routine_id AND routine.deleted_at IS NULL`).
			Joins(`INNER JOIN "UsersToStationsTable" uts ON uts.station_id = routine.station_id`).
			Where(`"RoutineTaskTable".id = "RoutineTaskRecordTable".routine_task_id`).
			Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions)
		return db.Where(`"RoutineTaskRecordTable".id IN ? AND EXISTS (?)`, ids, subQuery)
	}
}

func (sc *RoutineTaskRecordScope) IncludePreloads(preloads []schemas.RoutineTaskRecordRelation) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
		return db
	}
}
