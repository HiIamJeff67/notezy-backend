package repositories

import (
	"gorm.io/gorm/clause"
	"net/http"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/options"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/scopes"
	durablejobexceptions "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/exceptions"
	partialupdate "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/partialupdate"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type UserAccountRepositoryInterface interface {
	GetOneByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UserAccount, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserAccountInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserAccountInput, opts ...options.RepositoryOptions) (*schemas.UserAccount, *exceptions.Exception)
}

type UserAccountRepository struct{}

func NewUserAccountRepository() UserAccountRepositoryInterface {
	return &UserAccountRepository{}
}

func (r *UserAccountRepository) GetOneByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UserAccount, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var userAccount schemas.UserAccount
	result := parsedOptions.DB.Table(schemas.UserAccount{}.TableName()).
		Where("user_id = ?", userId).
		Scopes(scopes.Locking(parsedOptions.LockingStrength)).
		First(&userAccount)
	if err := result.Error; err != nil {
		return nil, durablejobexceptions.UserAccount.NotFound().WithOrigin(err)
	}

	return &userAccount, nil
}

func (r *UserAccountRepository) CreateOneByUserId(
	userId uuid.UUID,
	input inputs.CreateUserAccountInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var newUserAccount schemas.UserAccount
	newUserAccount.UserId = userId

	if err := copier.Copy(&newUserAccount, &input); err != nil {
		return nil, durablejobexceptions.UserAccount.FailedToCreate().WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.UserAccount{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUserAccount)
	if err := result.Error; err != nil {
		return nil, durablejobexceptions.UserAccount.FailedToCreate().WithOrigin(err)
	}

	return &newUserAccount.Id, nil
}

func (r *UserAccountRepository) UpdateOneByUserId(
	userId uuid.UUID,
	input inputs.PartialUpdateUserAccountInput,
	opts ...options.RepositoryOptions,
) (*schemas.UserAccount, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	existingUserAccount, exception := r.GetOneByUserId(
		userId,
		opts...,
	)
	if exception = exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
		{First: existingUserAccount == nil, Second: durablejobexceptions.UserAccount.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	updates, err := partialupdate.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserAccount)
	if err != nil {
		return nil, exceptions.New("FailedToPreprocessPartialUpdate", "Repository", "Update", "Failed to preprocess partial update", http.StatusInternalServerError, true)
	}

	result := parsedOptions.DB.Model(&schemas.UserAccount{}).
		Where("user_id = ?", userId).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, durablejobexceptions.UserAccount.FailedToUpdate().WithOrigin(err)
	}
	if result.RowsAffected == 0 {
		return nil, durablejobexceptions.UserAccount.NoChanges()
	}

	return &updates, nil
}

// We do not allow to just delete the userAccount,
// instead, the userAccount is only deleted by deleting the user
// func DeleteUserAccount(userId uuid.UUID) (deletedUserAccount User, err error) {}
