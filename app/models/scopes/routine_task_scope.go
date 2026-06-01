package scopes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
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
		// Use gorm.DB.Session to build a fresh statement for the subquery to avoid inheriting outer query clauses (especially in UPDATE/DELETE).
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToStations{}).
			Select("1").
			Where("station_id = \"RoutineTaskTable\".station_id AND user_id = ? AND permission IN ?", userId, permissions)
		return db.Where("\"RoutineTaskTable\".id = ? AND EXISTS (?)", id, subQuery)
	}
}

func (sc *RoutineTaskScope) PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Use gorm.DB.Session to build a fresh statement for the subquery to avoid inheriting outer query clauses (especially in UPDATE/DELETE).
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToStations{}).
			Select("1").
			Where("station_id = \"RoutineTaskTable\".station_id AND user_id = ? AND permission IN ?", userId, permissions)
		return db.Where("\"RoutineTaskTable\".id IN ? AND EXISTS (?)", ids, subQuery)
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
