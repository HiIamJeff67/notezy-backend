package repositories

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"

	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	options "notezy-backend/app/options"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

type SubShelfRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.SubShelf, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.SubShelf, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation, opts ...options.RepositoryOptions) (*schemas.SubShelf, *exceptions.Exception)
	GetAllByRootShelfId(rootShelfId uuid.UUID, userId uuid.UUID, preloads []schemas.SubShelfRelation, opts ...options.RepositoryOptions) ([]schemas.SubShelf, *exceptions.Exception)
	CreateOneByRootShelfId(rootShelfId uuid.UUID, userId uuid.UUID, input inputs.CreateSubShelfInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	BulkCreateManyByRootShelfIds(userId uuid.UUID, input []inputs.BulkCreateSubShelfInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateSubShelfInput, opts ...options.RepositoryOptions) (*schemas.SubShelf, *exceptions.Exception)
	BulkUpdateManyByIds(userId uuid.UUID, input []inputs.BulkUpdateSubShelfInput, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.SubShelf, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.SubShelf, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type SubShelfRepository struct{}

func NewSubShelfRepository() SubShelfRepositoryInterface {
	return &SubShelfRepository{}
}

func (r *SubShelfRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("id = ? AND EXISTS (?)", id, subQuery)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *SubShelfRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	query := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("id IN ? AND EXISTS (?)", ids, subQuery)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("deleted_at IS NULL")
	}

	var count int64 = 0
	result := query.Count(&count)
	if err := result.Error; err != nil {
		return false
	}

	return count > 0
}

func (r *SubShelfRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.SubShelf, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subShelf := schemas.SubShelf{}

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("id = ? AND EXISTS (?)", id, subQuery)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"SubShelfTable\".deleted_at IS NOT NULL")
	case types.Ternary_Neutral:
		break
	case types.Ternary_Negative:
		query = query.Where("\"SubShelfTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.First(&subShelf)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.NotFound().WithOrigin(result.Error)},
		{First: subShelf.Id == uuid.Nil, Second: exceptions.Shelf.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &subShelf, nil
}

func (r *SubShelfRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.SubShelf, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subShelves := []schemas.SubShelf{}

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("id IN ? AND EXISTS (?)", ids, subQuery)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("\"SubShelfTable\".deleted_at IS NOT NULL")
	case types.Ternary_Neutral:
		break
	case types.Ternary_Negative:
		query = query.Where("\"SubShelfTable\".deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.Find(&subShelves)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.NotFound().WithOrigin(result.Error)},
		{First: len(subShelves) == 0, Second: exceptions.Shelf.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return subShelves, nil
}

func (r *SubShelfRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
	opts ...options.RepositoryOptions,
) (*schemas.SubShelf, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	return r.CheckPermissionAndGetOneById(
		id,
		userId,
		preloads,
		allowedPermissions,
		opts...,
	)
}

func (r *SubShelfRepository) GetAllByRootShelfId(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.SubShelfRelation,
	opts ...options.RepositoryOptions,
) ([]schemas.SubShelf, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Read,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Admin,
	}

	subShelves := []schemas.SubShelf{}

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("root_shelf_id = ? AND EXISTS (?)", rootShelfId, subQuery)
	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	result := query.Find(&subShelves)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.NotFound().WithOrigin(result.Error)},
		{First: len(subShelves) == 0, Second: exceptions.Shelf.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return subShelves, nil
}

func (r *SubShelfRepository) CreateOneByRootShelfId(
	rootShelfId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateSubShelfInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	var newSubShelf schemas.SubShelf
	if input.PrevSubShelfId != nil {
		prevSubShelf, exception := r.CheckPermissionAndGetOneById(
			*input.PrevSubShelfId,
			userId,
			nil,
			allowedPermissions,
			opts...,
		)
		if exception = exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
			{First: prevSubShelf.RootShelfId != rootShelfId, Second: exceptions.Shelf.InvalidDto("the given prev sub shelf is not one of the children of the given root shelf")},
		}); exception != nil {
			return nil, exception
		}
		prevSubShelf.Path = append(prevSubShelf.Path, prevSubShelf.Id)
		newSubShelf.Path = prevSubShelf.Path
	} else {
		rootShelfRepository := NewRootShelfRepository()
		if !rootShelfRepository.HasPermission(
			rootShelfId,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.Shelf.NoPermission("create sub shelf by the given root shelf")
		}
	}

	if err := copier.Copy(&newSubShelf, &input); err != nil {
		return nil, exceptions.Shelf.InvalidInput().WithOrigin(err)
	}
	newSubShelf.RootShelfId = rootShelfId

	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newSubShelf)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToCreate().WithOrigin(result.Error)},
		{First: newSubShelf.Id == uuid.Nil, Second: exceptions.Shelf.FailedToCreate()},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &newSubShelf.Id, nil
}

func (r *SubShelfRepository) BulkCreateManyByRootShelfIds(
	userId uuid.UUID,
	input []inputs.BulkCreateSubShelfInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(input) == 0 {
		return nil, exceptions.Shelf.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	isPrevSubShelfExist := make(map[uuid.UUID]bool)
	isRootShelfExist := make(map[uuid.UUID]bool)
	prevSubShelfIds := make([]uuid.UUID, len(input))
	rootShelfIds := make([]uuid.UUID, len(input))
	for index, in := range input {
		if in.PrevSubShelfId != nil {
			if isPrevSubShelfExist[*in.PrevSubShelfId] {
				prevSubShelfIds[index] = *in.PrevSubShelfId
			}
			isPrevSubShelfExist[*in.PrevSubShelfId] = true
		}
		if isRootShelfExist[in.RootShelfId] {
			rootShelfIds[index] = in.RootShelfId
		}
		isRootShelfExist[in.RootShelfId] = true
	}

	validPrevSubShelves, exception := r.CheckPermissionsAndGetManyByIds(
		prevSubShelfIds,
		userId,
		nil,
		allowedPermissions,
	)
	if exception != nil {
		return nil, exception
	}
	isPrevSubShelfValid := make(map[uuid.UUID]*uuid.UUID)
	for _, validPrevSubShelf := range validPrevSubShelves {
		isPrevSubShelfValid[validPrevSubShelf.Id] = &validPrevSubShelf.RootShelfId
	}

	rootShelfRepository := NewRootShelfRepository()
	validRootShelves, exception := rootShelfRepository.CheckPermissionsAndGetManyByIds(
		rootShelfIds,
		userId,
		nil,
		allowedPermissions,
	)
	if exception != nil {
		return nil, exception
	}
	isRootShelfValid := make(map[uuid.UUID]bool)
	for _, validRootShelf := range validRootShelves {
		isRootShelfValid[validRootShelf.Id] = true
	}

	var newSubShelves []schemas.SubShelf
	for _, in := range input {
		if !isRootShelfValid[in.RootShelfId] ||
			(in.PrevSubShelfId != nil && (isPrevSubShelfValid[*in.PrevSubShelfId] != &in.RootShelfId)) {
			continue
		}
		var newSubShelf schemas.SubShelf
		if err := copier.Copy(&newSubShelf, &in); err != nil {
			return nil, exceptions.Shelf.InvalidInput().WithOrigin(err)
		}
		newSubShelf.RootShelfId = in.RootShelfId
		newSubShelves = append(newSubShelves, newSubShelf)
	}

	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		CreateInBatches(&newSubShelves, parsedOptions.BatchSize)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToCreate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	newSubShelfIds := make([]uuid.UUID, len(newSubShelves))
	for index, newSubShelf := range newSubShelves {
		newSubShelfIds[index] = newSubShelf.Id
	}

	return newSubShelfIds, nil
}

func (r *SubShelfRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateSubShelfInput,
	opts ...options.RepositoryOptions,
) (*schemas.SubShelf, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	existingSubShelf, exception := r.CheckPermissionAndGetOneById(
		id,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingSubShelf)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingSubShelf).WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &updates, nil
}

func (r *SubShelfRepository) BulkUpdateManyByIds(
	userId uuid.UUID,
	input []inputs.BulkUpdateSubShelfInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	isSubShelfValid := make(map[uuid.UUID]bool)
	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}
		subShelfIds := make([]uuid.UUID, len(input))
		for index, in := range input {
			subShelfIds[index] = in.Id
		}

		validSubShelves, exception := r.CheckPermissionsAndGetManyByIds(
			subShelfIds,
			userId,
			nil,
			allowedPermissions,
			opts...,
		)
		if exception != nil {
			return exception
		}

		for _, validSubShelf := range validSubShelves {
			isSubShelfValid[validSubShelf.Id] = true
		}
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, in := range input {
		if !parsedOptions.SkipPermissionCheck && !isSubShelfValid[in.Id] {
			continue
		}

		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::string)")
		valueArgs = append(valueArgs,
			in.Id,
			in.PartialUpdateInput.Values.Name,
		)
	}

	sql := fmt.Sprintf(`
		UPDATE "SubShelfTable" AS s
		SET
			name = COALESCE(v.name::string, s.name)
		FROM (VALUES %s) AS v(id, name)
		WHERE s.id = v.id::uuid AND s.deleted_at IS NULL
	`, strings.Join(valuePlaceholders, ","))
	result := parsedOptions.DB.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *SubShelfRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.SubShelf, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredSubShelf schemas.SubShelf
	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := parsedOptions.DB.Model(&restoredSubShelf).
		Clauses(clause.Returning{}).
		Where("id = ? AND EXISTS (?) AND deleted_at IS NOT NULL", id, subQuery).
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: restoredSubShelf.Id == uuid.Nil, Second: exceptions.Shelf.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &restoredSubShelf, nil
}

func (r *SubShelfRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.SubShelf, *exceptions.Exception) {
	if len(ids) == 0 {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredSubShelves []schemas.SubShelf
	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := parsedOptions.DB.Model(restoredSubShelves).
		Clauses(clause.Returning{}).
		Where("id IN ? AND EXISTS (?) AND deleted_at IS NOT NULL", ids, subQuery).
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: len(restoredSubShelves) == 0, Second: exceptions.Shelf.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return restoredSubShelves, nil
}

func (r *SubShelfRepository) SoftDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("id = ? AND EXISTS (?) AND deleted_at IS NULL", id, subQuery).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *SubShelfRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("id IN ? AND EXISTS (?) AND deleted_at IS NULL", ids, subQuery).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *SubShelfRepository) HardDeleteOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("id = ? AND EXISTS (?) AND deleted_at IS NOT NULL", id, subQuery).
		Delete(&schemas.SubShelf{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *SubShelfRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.BlockGroup.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id AND user_id = ? AND permission IN ?", userId, allowedPermissions)
	result := parsedOptions.DB.Model(&schemas.SubShelf{}).
		Where("id IN ? AND EXISTS (?) AND deleted_at IS NOT NULL", ids, subQuery).
		Delete(&schemas.SubShelf{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Shelf.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Shelf.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}
