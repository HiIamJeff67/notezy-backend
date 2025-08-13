package repositories

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	exceptions "notezy-backend/app/exceptions"
	models "notezy-backend/app/models"
	inputs "notezy-backend/app/models/inputs"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
)

/* ============================== Definitions ============================== */

type ShelfRepositoryInterface interface {
	GetOneById(id uuid.UUID, ownerId uuid.UUID, preloads *[]schemas.ShelfRelations) (*schemas.Shelf, *exceptions.Exception)
	GetOneByName(name string, ownerId uuid.UUID, preloads *[]schemas.ShelfRelations) (*schemas.Shelf, *exceptions.Exception)
	CreateOneByOwnerId(ownerId uuid.UUID, input inputs.CreateShelfInput) (*uuid.UUID, *exceptions.Exception)
	UpdateOneById(id uuid.UUID, ownerId uuid.UUID, input inputs.PartialUpdateShelfInput) (*schemas.Shelf, *exceptions.Exception)
	DeleteOneById(id uuid.UUID, ownerId uuid.UUID) *exceptions.Exception
}

type ShelfRepository struct {
	db *gorm.DB
}

func NewShelfRepository(db *gorm.DB) ShelfRepositoryInterface {
	if db == nil {
		db = models.NotezyDB
	}
	return &ShelfRepository{db: db}
}

/* ============================== CRUD operations ============================== */

func (r *ShelfRepository) GetOneById(id uuid.UUID, ownerId uuid.UUID, preloads *[]schemas.ShelfRelations) (*schemas.Shelf, *exceptions.Exception) {
	shelf := schemas.Shelf{}
	db := r.db.Table(schemas.Shelf{}.TableName())
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("id = ? AND owner_id = ?", id, ownerId).
		First(&shelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &shelf, nil
}

func (r *ShelfRepository) GetOneByName(name string, ownerId uuid.UUID, preloads *[]schemas.ShelfRelations) (*schemas.Shelf, *exceptions.Exception) {
	shelf := schemas.Shelf{}
	db := r.db.Table(schemas.Shelf{}.TableName())
	if preloads != nil {
		for _, preload := range *preloads {
			db = db.Preload(string(preload))
		}
	}

	result := db.Where("name = ? AND owner_id = ?", name, ownerId).
		First(&shelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.NotFound().WithError(err)
	}

	return &shelf, nil
}

func (r *ShelfRepository) CreateOneByOwnerId(ownerId uuid.UUID, input inputs.CreateShelfInput) (*uuid.UUID, *exceptions.Exception) {
	var newShelf schemas.Shelf
	newShelf.OwnerId = ownerId
	rootNode, exception := util.NewShelfNode(ownerId, input.Name, nil)
	if exception != nil {
		return nil, exception
	}
	encodedStructure, exception := util.EncodeShelfNode(rootNode)
	if exception != nil {
		return nil, exception
	}
	newShelf.EncodedStructure = encodedStructure
	if err := copier.Copy(&newShelf, &input); err != nil {
		return nil, exceptions.Theme.FailedToCreate().WithError(err)
	}

	result := r.db.Table(schemas.Shelf{}.TableName()).
		Create(&newShelf)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToCreate().WithError(err)
	}

	return &newShelf.Id, nil
}

func (r *ShelfRepository) UpdateOneById(id uuid.UUID, ownerId uuid.UUID, input inputs.PartialUpdateShelfInput) (*schemas.Shelf, *exceptions.Exception) {
	existingShelf, exception := r.GetOneById(id, ownerId, nil)
	if exception != nil || existingShelf == nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingShelf)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingShelf)
	}

	result := r.db.Table(schemas.Shelf{}.TableName()).
		Where("id = ? AND owner_id = ?", id, ownerId).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.Shelf.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 { // check if we do update it or not
		return nil, exceptions.Shelf.NoChanges()
	}

	return &updates, nil
}

func (r *ShelfRepository) DeleteOneById(id uuid.UUID, ownerId uuid.UUID) *exceptions.Exception {
	var deletedShelf schemas.Shelf

	result := r.db.Table(schemas.Shelf{}.TableName()).
		Where("id = ? AND owner_id = ?", id, ownerId).
		Clauses(clause.Returning{}).
		Delete(&deletedShelf)
	if err := result.Error; err != nil {
		return exceptions.Shelf.FailedToDelete().WithError(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.Shelf.NotFound()
	}

	return nil
}
