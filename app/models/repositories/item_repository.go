package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	array "github.com/HiIamJeff67/notezy-backend/shared/lib/array"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type ItemRepositoryInterface interface {
	HasPermission(id uuid.UUID, itemType enums.ItemType, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(itemIdentities []types.Pair[uuid.UUID, enums.ItemType], userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, itemType enums.ItemType, userId uuid.UUID, preloads []schemas.ItemRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.Item, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(itemIdentities []types.Pair[uuid.UUID, enums.ItemType], userId uuid.UUID, preloads []schemas.ItemRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.Item, *exceptions.Exception)
}

type ItemRepository struct {
	itemScope scopes.ItemScopeInterface
}

func NewItemRepository(
	itemScope scopes.ItemScopeInterface,
) ItemRepositoryInterface {
	return &ItemRepository{
		itemScope: itemScope,
	}
}

func (r *ItemRepository) HasPermission(
	id uuid.UUID,
	itemType enums.ItemType,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.Item{}).
		Select("1").
		Scopes(r.itemScope.PassPermissionCheck(id, itemType, userId, allowedPermissions)).
		Scopes(r.itemScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Limit(1).
		Scan(&marker)
	if err := result.Error; err != nil {
		return false
	}

	return marker == 1
}

func (r *ItemRepository) HavePermissions(
	itemIdentities []types.Pair[uuid.UUID, enums.ItemType],
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedItems []schemas.Item
	result := parsedOptions.DB.
		Model(&schemas.Item{}).
		Select(`DISTINCT "ItemTable".id, "ItemTable".type`).
		Scopes(r.itemScope.PassPermissionChecks(itemIdentities, userId, allowedPermissions)).
		Scopes(r.itemScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&permittedItems)
	if err := result.Error; err != nil {
		return false
	}

	permittedItemIdentities := make([]types.Pair[uuid.UUID, enums.ItemType], len(permittedItems))
	for index, permittedItem := range permittedItems {
		permittedItemIdentities[index] = types.Pair[uuid.UUID, enums.ItemType]{
			First:  permittedItem.Id,
			Second: permittedItem.Type,
		}
	}

	return array.GetDistinctCount(itemIdentities) == array.GetDistinctCount(permittedItemIdentities)
}

func (r *ItemRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	itemType enums.ItemType,
	userId uuid.UUID,
	preloads []schemas.ItemRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.Item, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var item schemas.Item
	result := parsedOptions.DB.
		Model(&schemas.Item{}).
		Scopes(r.itemScope.PassPermissionCheck(id, itemType, userId, allowedPermissions)).
		Scopes(r.itemScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.itemScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&item)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Item.NotFound().WithOrigin(result.Error)},
		{First: item.Id == uuid.Nil, Second: exceptions.Item.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &item, nil
}

func (r *ItemRepository) CheckPermissionsAndGetManyByIds(
	itemIdentities []types.Pair[uuid.UUID, enums.ItemType],
	userId uuid.UUID,
	preloads []schemas.ItemRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.Item, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var items []schemas.Item
	result := parsedOptions.DB.
		Model(&schemas.Item{}).
		Scopes(r.itemScope.PassPermissionChecks(itemIdentities, userId, allowedPermissions)).
		Scopes(r.itemScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.itemScope.IncludePreloads(preloads)).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&items)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Item.NotFound().WithOrigin(result.Error)},
		{First: len(items) == 0, Second: exceptions.Item.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return items, nil
}
