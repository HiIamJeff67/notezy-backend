package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	array "github.com/HiIamJeff67/notezy-backend/shared/lib/array"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type MaterialRepositoryInterface interface {
	HasPermission(id uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	HavePermissions(ids []uuid.UUID, userId uuid.UUID, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) bool
	CheckPermissionAndGetOneById(id uuid.UUID, userId uuid.UUID, preloads []schemas.MaterialRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) (*schemas.Material, *exceptions.Exception)
	CheckPermissionsAndGetManyByIds(ids []uuid.UUID, userId uuid.UUID, preloads []schemas.MaterialRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]schemas.Material, *exceptions.Exception)
	GetOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.Material, *exceptions.Exception)
	CreateOneBySubShelfId(subShelfId uuid.UUID, userId uuid.UUID, input inputs.CreateMaterialInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, userId uuid.UUID, input inputs.PartialUpdateMaterialInput, opts ...options.RepositoryOptions) (*schemas.Material, *exceptions.Exception)
	RestoreSoftDeletedOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.Material, *exceptions.Exception)
	RestoreSoftDeletedManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.Material, *exceptions.Exception)
	SoftDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	SoftDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteOneById(id uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
	HardDeleteManyByIds(ids []uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception

	/* ============================== System Only Method ============================== */

	BulkCheckPermissionsAndGetManyByIds(inputs []inputs.BulkCheckMaterialPermissionInput, preloads []schemas.MaterialRelation, allowedPermissions []enums.AccessControlPermission, opts ...options.RepositoryOptions) ([]bool, []schemas.Material, *exceptions.Exception)
	BulkDeleteMany(inputs []inputs.BulkDeleteMaterialInput, opts ...options.RepositoryOptions) ([]bool, *exceptions.Exception)
}

type MaterialRepository struct {
	materialScope scopes.MaterialScopeInterface
}

func NewMaterialRepository(materialScope scopes.MaterialScopeInterface) MaterialRepositoryInterface {
	return &MaterialRepository{
		materialScope: materialScope,
	}
}

func (r *MaterialRepository) HasPermission(
	id uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var marker int
	result := parsedOptions.DB.
		Model(&schemas.Material{}).
		Select("1").
		Scopes(r.materialScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Limit(1).
		Scan(&marker)
	if err := result.Error; err != nil {
		return false
	}

	return marker == 1
}

func (r *MaterialRepository) HavePermissions(
	ids []uuid.UUID,
	userId uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) bool {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var permittedIds []uuid.UUID
	result := parsedOptions.DB.
		Model(&schemas.Material{}).
		Select(`DISTINCT "MaterialTable".id`).
		Scopes(r.materialScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&permittedIds)
	if err := result.Error; err != nil {
		return false
	}

	return array.GetDistinctCount(ids) == array.GetDistinctCount(permittedIds)
}

func (r *MaterialRepository) CheckPermissionAndGetOneById(
	id uuid.UUID,
	userId uuid.UUID,
	preloads []schemas.MaterialRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) (*schemas.Material, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var material schemas.Material
	result := parsedOptions.DB.
		Model(&schemas.Material{}).
		Scopes(r.materialScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.materialScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&material)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Material.NotFound().WithOrigin(result.Error)},
		{First: material.Id == uuid.Nil, Second: exceptions.Material.NotFound()},
	}); exception != nil {
		return nil, exception
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

	var materials []schemas.Material
	result := parsedOptions.DB.
		Model(&schemas.Material{}).
		Scopes(r.materialScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.materialScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&materials)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Material.NotFound().WithOrigin(result.Error)},
		{First: len(materials) == 0, Second: exceptions.Material.NotFound()},
	}); exception != nil {
		return nil, exception
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

	shouldStartTransaction := !parsedOptions.IsTransactionStarted && !parsedOptions.SkipPermissionCheck
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	}

	if !parsedOptions.SkipPermissionCheck {
		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		subShelfRepository := NewSubShelfRepository(scopes.NewSubShelfScope())

		if !subShelfRepository.HasPermission(
			subShelfId,
			userId,
			allowedPermissions,
			opts...,
		) {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Shelf.NoPermission("create a material under this shelf")
		}
	}

	var newMaterial schemas.Material
	if err := copier.Copy(&newMaterial, &input); err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Material.FailedToCreate().WithOrigin(err)
	}
	newMaterial.ParentSubShelfId = subShelfId

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Create(&newMaterial)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Material.FailedToCreate().WithOrigin(result.Error)},
		{First: newMaterial.Id == uuid.Nil, Second: exceptions.Material.FailedToCreate()},
		{First: result.RowsAffected == 0, Second: exceptions.Material.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Material.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &newMaterial.Id, nil
}

func (r *MaterialRepository) UpdateOneById(
	id uuid.UUID,
	userId uuid.UUID,
	input inputs.PartialUpdateMaterialInput,
	opts ...options.RepositoryOptions,
) (*schemas.Material, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
		opts = append(opts, options.WithTransactionDB(parsedOptions.DB))
		opts = append(opts, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	}

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
		parsedOptions.DB.Rollback()
		return nil, exception
	}
	if existingMaterial == nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Material.NotFound()
	}

	// if the root shelf id is required to be updated in the database
	if input.Values.ParentSubShelfId != nil && !util.CheckSetNull(input.SetNull, "ParentSubShelfId") {
		subShelfRepository := NewSubShelfRepository(scopes.NewSubShelfScope())
		// check if the user has the enough permission to the destination shelf
		if !subShelfRepository.HasPermission(
			*input.Values.ParentSubShelfId,
			userId,
			allowedPermissions,
			opts...,
		) {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Shelf.NoPermission("move a material to this shelf")
		}
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingMaterial)
	if err != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
			input.Values,
			input.SetNull,
			*existingMaterial,
		).WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Where("id = ? AND deleted_at IS NULL", id). // no need to check the permission here, since we have done that part on the above
		Select("*").
		Updates(&updates)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Material.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Material.NoChanges()},
	}); exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Material.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	return &updates, nil
}

func (r *MaterialRepository) RestoreSoftDeletedOneById(
	id uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.Material, *exceptions.Exception) {
	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredMaterial schemas.Material
	query := parsedOptions.DB.Model(&restoredMaterial).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted))
	if !parsedOptions.SkipPermissionCheck {
		query = query.Scopes(r.materialScope.PassPermissionCheck(id, userId, allowedPermissions))
	}

	result := query.
		Clauses(clause.Returning{}).
		Where(`"MaterialTable".id = ?`, id).
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Material.FailedToUpdate().WithOrigin(result.Error)},
		{First: restoredMaterial.Id == uuid.Nil, Second: exceptions.Material.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Material.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &restoredMaterial, nil
}

func (r *MaterialRepository) RestoreSoftDeletedManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.Material, *exceptions.Exception) {
	if len(ids) == 0 {
		return nil, exceptions.Material.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	var restoredMaterials []schemas.Material
	query := parsedOptions.DB.Model(&restoredMaterials).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted))
	if !parsedOptions.SkipPermissionCheck {
		query = query.Scopes(r.materialScope.PassPermissionChecks(ids, userId, allowedPermissions))
	}

	result := query.
		Clauses(clause.Returning{}).
		Where(`"MaterialTable".id IN ?`, ids).
		Updates(map[string]interface{}{"deleted_at": nil}) // force to assign null value
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Material.FailedToUpdate().WithOrigin(result.Error)},
		{First: len(restoredMaterials) != len(ids), Second: exceptions.Material.FailedToUpdate()},
		{First: result.RowsAffected == 0, Second: exceptions.Material.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return restoredMaterials, nil
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

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Scopes(r.materialScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"MaterialTable".id = ?`, id).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Material.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Material.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *MaterialRepository) SoftDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.Material.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Negative))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Scopes(r.materialScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"MaterialTable".id IN ?`, ids).
		Update("deleted_at", time.Now())
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Material.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Material.NoChanges()},
	}); exception != nil {
		return exception
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

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Scopes(r.materialScope.PassPermissionCheck(id, userId, allowedPermissions)).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"MaterialTable".id = ?`, id).
		Delete(&schemas.Material{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Material.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Material.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

func (r *MaterialRepository) HardDeleteManyByIds(
	ids []uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(ids) == 0 {
		return exceptions.Material.NoChanges()
	}

	opts = append(opts, options.WithOnlyDeleted(types.Ternary_Positive))
	parsedOptions := options.ParseRepositoryOptions(opts...)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}

	result := parsedOptions.DB.Model(&schemas.Material{}).
		Scopes(r.materialScope.PassPermissionChecks(ids, userId, allowedPermissions)).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Where(`"MaterialTable".id IN ?`, ids).
		Delete(&schemas.Material{})
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Material.FailedToDelete().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.Material.NoChanges()},
	}); exception != nil {
		return exception
	}

	return nil
}

/* ============================== System Only Method ============================== */

func (r *MaterialRepository) BulkCheckPermissionsAndGetManyByIds(
	inputs []inputs.BulkCheckMaterialPermissionInput,
	preloads []schemas.MaterialRelation,
	allowedPermissions []enums.AccessControlPermission,
	opts ...options.RepositoryOptions,
) ([]bool, []schemas.Material, *exceptions.Exception) {
	if len(inputs) == 0 {
		return []bool{}, []schemas.Material{}, nil
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	successes := make([]bool, len(inputs))
	ids := make([]uuid.UUID, 0, len(inputs))
	userIds := make([]uuid.UUID, 0, len(inputs))
	for _, in := range inputs {
		ids = append(ids, in.Id)
		userIds = append(userIds, in.UserId)
	}

	var validTargets []struct {
		Id     uuid.UUID `gorm:"column:id"`
		UserId uuid.UUID `gorm:"column:user_id"`
	}
	result := parsedOptions.DB.Model(&schemas.Material{}).
		Select(`"MaterialTable".id, uts.user_id`).
		Joins(`INNER JOIN "SubShelfTable" AS ss ON ss.id = "MaterialTable".parent_sub_shelf_id`).
		Joins(`INNER JOIN "UsersToShelvesTable" AS uts ON uts.root_shelf_id = ss.root_shelf_id`).
		Where(`"MaterialTable".id IN ?`, ids).
		Where("uts.user_id IN ? AND uts.permission IN ?", userIds, allowedPermissions).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scan(&validTargets)
	if result.Error != nil {
		return nil, nil, exceptions.Material.NotFound().WithOrigin(result.Error)
	}

	validTargetByUserId := make(map[[2]uuid.UUID]bool, len(validTargets))
	for _, validTarget := range validTargets {
		validTargetByUserId[[2]uuid.UUID{validTarget.Id, validTarget.UserId}] = true
	}

	validIdSet := make(map[uuid.UUID]bool, len(validTargets))
	for _, in := range inputs {
		if validTargetByUserId[[2]uuid.UUID{in.Id, in.UserId}] {
			validIdSet[in.Id] = true
		}
	}

	validIds := make([]uuid.UUID, 0, len(validIdSet))
	for validId := range validIdSet {
		validIds = append(validIds, validId)
	}
	if len(validIds) == 0 {
		return successes, []schemas.Material{}, nil
	}

	var materials []schemas.Material
	result = parsedOptions.DB.Model(&schemas.Material{}).
		Where(`"MaterialTable".id IN ?`, validIds).
		Scopes(r.materialScope.FilterOnlyDeleted(parsedOptions.OnlyDeleted)).
		Scopes(r.materialScope.IncludePreloads(preloads)).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		Find(&materials)
	if result.Error != nil {
		return nil, nil, exceptions.Material.NotFound().WithOrigin(result.Error)
	}

	foundIdSet := make(map[uuid.UUID]bool, len(materials))
	for _, material := range materials {
		foundIdSet[material.Id] = true
	}
	for index, in := range inputs {
		if validTargetByUserId[[2]uuid.UUID{in.Id, in.UserId}] && foundIdSet[in.Id] {
			successes[index] = true
		}
	}

	return successes, materials, nil
}

func (r *MaterialRepository) BulkDeleteMany(
	bulkInputs []inputs.BulkDeleteMaterialInput,
	opts ...options.RepositoryOptions,
) ([]bool, *exceptions.Exception) {
	if len(bulkInputs) == 0 {
		return []bool{}, exceptions.Material.NoChanges()
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	shouldStartTransaction := !parsedOptions.IsTransactionStarted
	if shouldStartTransaction {
		parsedOptions.DB = parsedOptions.DB.Begin()
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	checkInputs := make([]inputs.BulkCheckMaterialPermissionInput, len(bulkInputs))
	for index, in := range bulkInputs {
		checkInputs[index] = inputs.BulkCheckMaterialPermissionInput{
			UserId: in.UserId,
			Id:     in.Id,
		}
	}
	checkOptions := append(opts, options.WithTransactionDB(parsedOptions.DB))
	checkOptions = append(checkOptions, options.WithOnlyDeleted(types.Ternary_Negative))
	checkOptions = append(checkOptions, options.WithLockingStrength(options.LockingStrengthNoKeyUpdate))
	successes, _, exception := r.BulkCheckPermissionsAndGetManyByIds(checkInputs, nil, allowedPermissions, checkOptions...)
	if exception != nil {
		parsedOptions.DB.Rollback()
		return nil, exception
	}

	validIds := make([]uuid.UUID, 0, len(bulkInputs))
	for index, in := range bulkInputs {
		if successes[index] {
			validIds = append(validIds, in.Id)
		}
	}
	if len(validIds) == 0 {
		if shouldStartTransaction {
			parsedOptions.DB.Rollback()
		}
		return successes, nil
	}

	var deletedMaterials []schemas.Material
	result := parsedOptions.DB.Model(&deletedMaterials).
		Clauses(clause.Returning{}).
		Where("id IN ? AND deleted_at IS NULL", validIds).
		Updates(map[string]interface{}{"deleted_at": time.Now(), "updated_at": time.Now()})
	if result.Error != nil {
		parsedOptions.DB.Rollback()
		return nil, exceptions.Material.FailedToDelete().WithOrigin(result.Error)
	}

	if shouldStartTransaction {
		if err := parsedOptions.DB.Commit().Error; err != nil {
			parsedOptions.DB.Rollback()
			return nil, exceptions.Material.FailedToCommitTransaction().WithOrigin(err)
		}
	}

	deletedIdSet := make(map[uuid.UUID]bool, len(deletedMaterials))
	for _, deletedMaterial := range deletedMaterials {
		deletedIdSet[deletedMaterial.Id] = true
	}
	for index, in := range bulkInputs {
		if successes[index] && deletedIdSet[in.Id] {
			successes[index] = true
		} else {
			successes[index] = false
		}
	}

	return successes, nil
}
