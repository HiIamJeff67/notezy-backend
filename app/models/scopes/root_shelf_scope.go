package scopes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RootShelfScopeInterface interface {
	PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	FilterOnlyDeleted(onlyDeleted types.Ternary) func(db *gorm.DB) *gorm.DB
	IncludePreloads(preloads []schemas.RootShelfRelation) func(db *gorm.DB) *gorm.DB
}

type RootShelfScope struct{}

func NewRootShelfScope() RootShelfScopeInterface {
	return &RootShelfScope{}
}

func (sc *RootShelfScope) PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		allowedPermissions := contexts.IntersectAllowedPermissions(db.Statement.Context, permissions)

		// Use gorm.DB.Session to build a fresh statement for the subquery to avoid inheriting outer query clauses (especially in UPDATE/DELETE).
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToShelves{}).
			Select("1").
			Where("root_shelf_id = \"RootShelfTable\".id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
		return db.Where("id = ? AND EXISTS (?)", id, subQuery)
	}
}

func (sc *RootShelfScope) PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		allowedPermissions := contexts.IntersectAllowedPermissions(db.Statement.Context, permissions)

		// Use gorm.DB.Session to build a fresh statement for the subquery to avoid inheriting outer query clauses (especially in UPDATE/DELETE).
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToShelves{}).
			Select("1").
			Where("root_shelf_id = \"RootShelfTable\".id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
		return db.Where("id IN ? AND EXISTS (?)", ids, subQuery)
	}
}

func (sc *RootShelfScope) FilterOnlyDeleted(onlyDeleted types.Ternary) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch onlyDeleted {
		case types.Ternary_Positive:
			return db.Where("\"RootShelfTable\".deleted_at IS NOT NULL")
		case types.Ternary_Negative:
			return db.Where("\"RootShelfTable\".deleted_at IS NULL")
		default:
			return db
		}
	}
}

func (sc *RootShelfScope) IncludePreloads(preloads []schemas.RootShelfRelation) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
		return db
	}
}
