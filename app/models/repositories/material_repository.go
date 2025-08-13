package repositories

import (
	"notezy-backend/app/exceptions"
	"notezy-backend/app/models"
	"notezy-backend/app/models/inputs"
	"notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/* ============================== Definitions ============================== */

type MaterialRepositoryInterface interface {
}

type MaterialRepository struct {
	db *gorm.DB
}

func NewMaterialRepository(db *gorm.DB) MaterialRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &MaterialRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *MaterialRepository) GetOneById(id uuid.UUID) (*schemas.Material, *exceptions.Exception) {
	material := schemas.Material{}

	result := r.db.Table(schemas.Material{}.TableName()).
		Where("id = ?", id).
		First(&material)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}

	return &material, nil
}

func (r *MaterialRepository) GetOneByName(name string) (*schemas.Material, *exceptions.Exception) {
	material := schemas.Material{}

	result := r.db.Table(schemas.Material{}.TableName()).
		Where("name = ?", name).
		First(&material)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.NotFound().WithError(err)
	}

	return &material, nil
}

func (r *MaterialRepository) CreateOne(input inputs.CreateMaterialInput) (*uuid.UUID, *exceptions.Exception) {
	var newMaterial schemas.Material
	if err := copier.Copy(&newMaterial, &input); err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}

	result := r.db.Table(schemas.Material{}.TableName()).
		Create(&newMaterial)
	if err := result.Error; err != nil {
		return nil, exceptions.Material.FailedToCreate().WithError(err)
	}

	return &newMaterial.Id, nil
}

func (r *MaterialRepository) UpdateOneById(id uuid.UUID, input inputs.PartialUpdateMaterialInput) (*schemas.Material, *exceptions.Exception) {
	existingMaterial, exception := r.GetOneById(id)
	if exception != nil || existingMaterial == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingMaterial)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingMaterial)
	}

	result := r.db.Table(schemas.Material{}.TableName()).
		Where("id = ?", id).
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

func (r *MaterialRepository) DeleteOneById(id uuid.UUID) *exceptions.Exception {
	var deletedMaterial schemas.Material

	result := r.db.Table(schemas.Material{}.TableName()).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Delete(&deletedMaterial)
	if err := result.Error; err != nil {
		return exceptions.Material.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Material.NotFound()
	}

	return nil
}
