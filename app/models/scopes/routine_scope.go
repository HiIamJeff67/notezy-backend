package scopes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutineScopeInterface interface {
	PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	FilterOnlyDeleted(onlyDeleted types.Ternary) func(db *gorm.DB) *gorm.DB
	IncludePreloads(preloads []schemas.RoutineRelation) func(db *gorm.DB) *gorm.DB
}

type RoutineScope struct{}

func NewRoutineScope() RoutineScopeInterface {
	return &RoutineScope{}
}

func (sc *RoutineScope) PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Use gorm.DB.Session to build a fresh statement for the subquery to avoid inheriting outer query clauses (especially in UPDATE/DELETE).
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToStations{}).
			Select("1").
			Where("station_id = \"RoutineTable\".station_id AND user_id = ? AND permission IN ?", userId, permissions)
		return db.Where("\"RoutineTable\".id = ? AND EXISTS (?)", id, subQuery)
	}
}

func (sc *RoutineScope) PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Use gorm.DB.Session to build a fresh statement for the subquery to avoid inheriting outer query clauses (especially in UPDATE/DELETE).
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToStations{}).
			Select("1").
			Where("station_id = \"RoutineTable\".station_id AND user_id = ? AND permission IN ?", userId, permissions)
		return db.Where("\"RoutineTable\".id IN ? AND EXISTS (?)", ids, subQuery)
	}
}

func (sc *RoutineScope) FilterOnlyDeleted(onlyDeleted types.Ternary) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch onlyDeleted {
		case types.Ternary_Positive:
			return db.Where("\"RoutineTable\".deleted_at IS NOT NULL")
		case types.Ternary_Negative:
			return db.Where("\"RoutineTable\".deleted_at IS NULL")
		default:
			return db
		}
	}
}

func (sc *RoutineScope) IncludePreloads(preloads []schemas.RoutineRelation) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
		return db
	}
}
