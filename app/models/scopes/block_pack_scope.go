package scopes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type BlockPackScopeInterface interface {
	PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	FilterOnlyDeleted(onlyDeleted types.Ternary) func(db *gorm.DB) *gorm.DB
	IncludePreloads(preloads []schemas.BlockPackRelation) func(db *gorm.DB) *gorm.DB
}

type BlockPackScope struct{}

func NewBlockPackScope() BlockPackScopeInterface {
	return &BlockPackScope{}
}

func (sc *BlockPackScope) PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Use gorm.DB.Session to build a fresh statement for the subquery to avoid inheriting outer query clauses (especially in UPDATE/DELETE).
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToShelves{}).
			Select("1").
			Where("root_shelf_id = ss.root_shelf_id").
			Where("user_id = ? AND permission IN ?", userId, permissions)
		return db.
			Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
			Where("\"BlockPackTable\".id = ? AND EXISTS (?)", id, subQuery)
	}
}

func (sc *BlockPackScope) PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Use gorm.DB.Session to build a fresh statement for the subquery to avoid inheriting outer query clauses (especially in UPDATE/DELETE).
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToShelves{}).
			Select("1").
			Where("root_shelf_id = ss.root_shelf_id").
			Where("user_id = ? AND permission IN ?", userId, permissions)
		return db.
			Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
			Where("\"BlockPackTable\".id IN ? AND EXISTS (?)", ids, subQuery)
	}
}

func (sc *BlockPackScope) FilterOnlyDeleted(onlyDeleted types.Ternary) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch onlyDeleted {
		case types.Ternary_Positive:
			return db.Where("\"BlockPackTable\".deleted_at IS NOT NULL")
		case types.Ternary_Negative:
			return db.Where("\"BlockPackTable\".deleted_at IS NULL")
		default:
			return db
		}
	}
}

func (sc *BlockPackScope) IncludePreloads(preloads []schemas.BlockPackRelation) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
		return db
	}
}
