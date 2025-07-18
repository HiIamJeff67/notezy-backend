package dataloaders

import (
	"context"
	"notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	"notezy-backend/app/services"
	"notezy-backend/shared/constants"

	gophersdataloader "github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

/* ============================== Enum Keys & Type ============================== */

type LoadUserSource string

const (
	LoadUserSource_ThemePublicId LoadUserSource = "LoadUserSourceThemePublicId"
)

type UserLoaderKey struct {
	PublicId string         `json:"publicId"`
	Source   LoadUserSource `json:"source"`
}

type UserLoaderType = gophersdataloader.Loader[UserLoaderKey, *gqlmodels.PublicUser]
type UserBatchFunctionType = gophersdataloader.BatchFunc[UserLoaderKey, *gqlmodels.PublicUser]
type UserResultType = gophersdataloader.Result[*gqlmodels.PublicUser]

/* ============================== Interface & Instance ============================== */

type UserDataloaderInterface interface {
	GetLoader() *UserLoaderType
	batchFunction() UserBatchFunctionType
}

type UserDataloader struct {
	userService services.UserServiceInterface
	loader      *UserLoaderType
}

func NewUserDataloader(db *gorm.DB) UserDataloaderInterface {
	dataloader := &UserDataloader{
		userService: services.NewUserService(db),
	}
	dataloader.loader = gophersdataloader.NewBatchedLoader(
		dataloader.batchFunction(),
		gophersdataloader.WithWait[UserLoaderKey, *gqlmodels.PublicUser](constants.LoaderDelayOfUser),
	)

	return dataloader
}

/* ============================== Dataloader Implementations ============================== */

func (d *UserDataloader) GetLoader() *UserLoaderType {
	return d.loader
}

func (d *UserDataloader) batchFunction() UserBatchFunctionType {
	return func(ctx context.Context, keys []UserLoaderKey) []*UserResultType {
		keysBySource := make(map[LoadUserSource][]string)
		keyToIndexesMap := make(map[UserLoaderKey][]int)

		for index, key := range keys {
			keysBySource[key.Source] = append(keysBySource[key.Source], key.PublicId)
			keyToIndexesMap[key] = append(keyToIndexesMap[key], index)
		}

		results := make([]*UserResultType, len(keys))

		for source, publicIds := range keysBySource {
			var publicUsers []*gqlmodels.PublicUser
			var exception *exceptions.Exception

			switch source {
			case LoadUserSource_ThemePublicId:
				publicUsers, exception = d.userService.GetPublicAuthorByThemePublicIds(ctx, publicIds)
			default:
				exception = exceptions.User.InvalidSourceInBatchFunction()
			}

			if exception != nil {
				for _, publicId := range publicIds {
					key := UserLoaderKey{PublicId: publicId, Source: source}
					if _, exists := keyToIndexesMap[key]; exists {
						for _, index := range keyToIndexesMap[key] {
							results[index] = &UserResultType{Error: exception.Error}
						}
					}
				}
				continue
			}

			for index, publicUser := range publicUsers {
				key := UserLoaderKey{PublicId: publicIds[index], Source: source}
				if _, exists := keyToIndexesMap[key]; exists {
					for _, originalIndex := range keyToIndexesMap[key] {
						results[originalIndex] = &UserResultType{Data: publicUser}
					}
				}
			}
		}

		return results
	}
}

/* ============================== Load Functions ============================== */

func (d *UserDataloader) LoadByThemePublicId(originalContext context.Context, publicId string) (*gqlmodels.PublicUser, error) {
	future := d.loader.Load(
		originalContext,
		UserLoaderKey{
			PublicId: publicId,
			Source:   LoadUserSource_ThemePublicId,
		},
	)

	publicUser, err := future()
	if err != nil {
		return nil, err
	}

	return publicUser, nil
}
