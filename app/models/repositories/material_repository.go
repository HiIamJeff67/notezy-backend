package repositories

import (
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

/* ============================== Definitions ============================== */

type MaterialRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HasPermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.MaterialRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.Material, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.MaterialRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.Material, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.Material, *exceptions.Exception)
	CreateOneBySubShelfId(subShelfId uuid.UUID, userId uuid.UUID, input inputs.CreateMaterialInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, matchedMaterialType *enums.MaterialType, input inputs.PartialUpdateMaterialInput, opts ...options.RepositoryOptions) (*schemas.Material, *exceptions.Exception)
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type MaterialRepository struct{}

func NewMaterialRepository() MaterialRepositoryInterface {
	return &MaterialRepository{}
}

/* ============================== Implementations ============================== */

func (r *MaterialRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.Material{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Where("id = ? AND EXISTS (?)",
			id, subQuery,
		)

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

func (r *MaterialRepository) HasPermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.Material{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Where("id IN ? EXISTS (?)",
			ids, subQuery,
		)

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

func (r *MaterialRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.MaterialRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.Material, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.Material{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Where("id = ? AND EXISTS (?)",
			id, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("deleted_at IS NOT NULL")
	case types.Ternary_Negative:
		query = query.Where("deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var material schemas.Material
	result := query.First(&material)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}

	return &material, nil
}

func (r *MaterialRepository) CheckPermissionsAndGetManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.MaterialRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]schemas.Material, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	subQuery := parsedOptions.DB.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			userId, allowedPermissions,
		)
	query := parsedOptions.DB.Model(&schemas.Material{}).
		Joins("INNER JOIN \"SubShelfTable\" ss ON parent_sub_shelf_id = ss.id").
		Where("id IN ? EXISTS (?)",
			ids, subQuery,
		)

	switch parsedOptions.OnlyDeleted {
	case types.Ternary_Positive:
		query = query.Where("deleted_at IS NOT NULL")
	case types.Ternary_Neutral:
		break
	case types.Ternary_Negative:
		query = query.Where("deleted_at IS NULL")
	}

	if len(preloads) > 0 {
		for _, preload := range preloads {
			query = query.Preload(string(preload))
		}
	}

	var materials []schemas.Material
	result := query.Find(&materials)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}
	if len(materials) == 0 {
		return nil, exceptions.Material.NotFound()
	}

	return materials, nil
}

func (r *MaterialRepository) GetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.Material, *exceptions.Exception) {
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	return r.CheckPermissionAndGetOneById(
		id,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
}

func (r *MaterialRepository) CreateOneBySubShelfId(
	subShelfId uuid.UUID,
	userId uuid.UUID,
	input inputs.CreateMaterialInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	subShelfRepository := NewSubShelfRepository()

	if !subShelfRepository.HasPermission(
		subShelfId,
		userId,
		allowedPermissions,
		opts...,
	) {
		return nil, exceptions.Shelf.NoPermission("create a material under this shelf")
	}

	var newMaterial schemas.Material
	if err := copier.Copy(&newMaterial, &input); err != nil {
		return nil, exceptions.Material.FailedToCreate().WithError(err)
	}
	newMaterial.ParentSubShelfId = subShelfId

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newMaterial)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.FailedToCreate().WithError(err)
	}

	return &newMaterial.Id, nil
}

func (r *MaterialRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	matchedMaterialType *enums.MaterialType,
	input inputs.PartialUpdateMaterialInput,
	opts ...options.RepositoryOptions,
) (*schemas.Material, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	// get and check the permission of the current user to the source shelf
	existingMaterial, exception := r.CheckPermissionAndGetOneById(
		id,
		userId,
		nil,
		allowedPermissions,
		opts...,
	)
	if exception != nil {
		return nil, exception
	}
	if existingMaterial == nil {
		return nil, exceptions.Material.NotFound()
	}

	// check if the material type is matched
	if matchedMaterialType != nil && existingMaterial.Type != *matchedMaterialType {
		return nil, exceptions.Material.MaterialTypeNotMatch(
			existingMaterial.Id.String(),
			existingMaterial.Type,
			matchedMaterialType,
		)
	}

	// if the root shelf id is required to be updated in the database
	if input.Values.ParentSubShelfId != nil && (input.SetNull == nil || !(*input.SetNull)["ParentSubShelfId"]) {
		subShelfRepository := NewSubShelfRepository()
		// check if the user has the enough permission to the destination shelf
		if !subShelfRepository.HasPermission(
			*input.Values.ParentSubShelfId,
			userId,
			allowedPermissions,
			opts...,
		) {
			return nil, exceptions.Shelf.NoPermission("move a material to this shelf")
		}
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingMaterial)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values,
			input.SetNull,
			*existingMaterial,
		).WithError(err)
	}

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Where("id = ? AND deleted_at IS NULL", id). // no need to check the permission here, since we have done that part on the above
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 { // check if we do update it or not
		return nil, exceptions.Material.NoChanges()
	}

	return &updates, nil
}

func (r *MaterialRepository) RestoreSoftDeletedOneById(
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

	if !r.HasPermission(
		id,
		userId,
		allowedPermissions,
		opts...,
	) {
		return exceptions.Material.NoPermission("restore a deleted material")
	}

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NoChanges()
	}

	return nil
}

func (r *MaterialRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	if !r.HasPermissions(
		ids,
		userId,
		allowedPermissions,
		opts...,
	) {
		return exceptions.Material.NoPermission("restore deleted materials")
	}

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Select("deleted_at").
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NoChanges()
	}

	return nil
}

func (r *MaterialRepository) SoftDeleteOneById(
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

	if !r.HasPermission(
		id,
		userId,
		allowedPermissions,
		opts...,
	) {
		return exceptions.Material.NoPermission("soft delete a material")
	}

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NoChanges()
	}

	return nil
}

func (r *MaterialRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	if !r.HasPermissions(
		ids,
		userId,
		allowedPermissions,
		opts...,
	) {
		return exceptions.Material.NoPermission("soft delete")
	}

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}

func (r *MaterialRepository) HardDeleteOneById(
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

	if !r.HasPermission(
		id,
		userId,
		allowedPermissions,
		opts...,
	) {
		return exceptions.Material.NoPermission("hard delete a material")
	}

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Delete(&schemas.Material{})
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NoChanges()
	}

	return nil
}

func (r *MaterialRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	if !r.HasPermissions(
		ids,
		userId,
		allowedPermissions,
		opts...,
	) {
		return exceptions.Material.NoPermission("hard delete")
	}

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Delete(&schemas.Material{})
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}
