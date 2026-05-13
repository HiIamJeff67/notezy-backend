package scopes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type BlockScopeInterface interface {
	PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB
	FilterOnlyDeleted(onlyDeleted types.Ternary) func(db *gorm.DB) *gorm.DB
	IncludePreloads(preloads []schemas.BlockRelation) func(db *gorm.DB) *gorm.DB
}

type BlockScope struct{}

func NewBlockScope() BlockScopeInterface {
	return &BlockScope{}
}

func (sc *BlockScope) PassPermissionCheck(id uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Use gorm.DB.Session to build a fresh statement for the subquery to avoid inheriting outer query clauses (especially in UPDATE/DELETE).
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToShelves{}).
			Select("1").
			Where("root_shelf_id = ss.root_shelf_id").
			Where("user_id = ? AND permission IN ?", userId, permissions)
		return db.
			Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
			Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
			Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
			Where("\"BlockTable\".id = ? AND EXISTS (?)", id, subQuery)
	}
}

func (sc *BlockScope) PassPermissionChecks(ids []uuid.UUID, userId uuid.UUID, permissions []enums.AccessControlPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Use gorm.DB.Session to build a fresh statement for the subquery to avoid inheriting outer query clauses (especially in UPDATE/DELETE).
		subQuery := db.Session(&gorm.Session{NewDB: true}).
			Model(&schemas.UsersToShelves{}).
			Select("1").
			Where("root_shelf_id = ss.root_shelf_id").
			Where("user_id = ? AND permission IN ?", userId, permissions)
		return db.
			Joins("INNER JOIN \"BlockGroupTable\" bg ON block_group_id = bg.id").
			Joins("INNER JOIN \"BlockPackTable\" bp ON bg.block_pack_id = bp.id").
			Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
			Where("\"BlockTable\".id IN ? AND EXISTS (?)", ids, subQuery)
	}
}

func (sc *BlockScope) FilterOnlyDeleted(onlyDeleted types.Ternary) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch onlyDeleted {
		case types.Ternary_Positive:
			return db.Where("\"BlockTable\".deleted_at IS NOT NULL")
		case types.Ternary_Negative:
			return db.Where("\"BlockTable\".deleted_at IS NULL")
		default:
			return db
		}
	}
}

func (sc *BlockScope) IncludePreloads(preloads []schemas.BlockRelation) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preloads {
			db = db.Preload(string(preload))
		}
		return db
	}
}
