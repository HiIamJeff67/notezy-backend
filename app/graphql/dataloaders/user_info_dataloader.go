package dataloaders

import (
	"context"

	gophersdataloader "github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"

	exceptions "notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

/* ============================== Enum Keys & Types ============================== */

type LoadUserInfoSource string

const (
	LoadUserInfoSource_UserPublicId LoadUserInfoSource = "LoadUserInfoSourceUserPublicId"
)

type UserInfoLoaderKey struct {
	PublicId string             `json:"publicId"`
	Source   LoadUserInfoSource `json:"source"`
}

type UserInfoLoaderType = gophersdataloader.Loader[UserInfoLoaderKey, *gqlmodels.PublicUserInfo]
type UserInfoBatchFunctionType = gophersdataloader.BatchFunc[UserInfoLoaderKey, *gqlmodels.PublicUserInfo]
type UserInfoResultType = gophersdataloader.Result[*gqlmodels.PublicUserInfo]

/* ============================== Interface & Instance ============================== */

type UserInfoDataloaderInterface interface {
	GetLoader() *UserInfoLoaderType
	batchFunction() UserInfoBatchFunctionType

	// load functions
	LoadByUserPublicId(originalContext context.Context, id string) (*gqlmodels.PublicUserInfo, error)
}

type UserInfoDataloader struct {
	userInfoService services.UserInfoServiceInterface
	loader          *UserInfoLoaderType
}

func NewUserInfoDataloader(db *gorm.DB) UserInfoDataloaderInterface {
	dataloader := &UserInfoDataloader{
		userInfoService: services.NewUserInfoService(db),
	}
	dataloader.loader = gophersdataloader.NewBatchedLoader(
		dataloader.batchFunction(),
		gophersdataloader.WithWait[UserInfoLoaderKey, *gqlmodels.PublicUserInfo](constants.LoaderDelayOfUserInfo),
	)

	return dataloader
}

/* ============================== Dataloader Implementations ============================== */

func (d *UserInfoDataloader) GetLoader() *UserInfoLoaderType {
	return d.loader
}

// this batch function will fetch the PublicUserInfos using the publicIds of the "PublicUsers"
func (d *UserInfoDataloader) batchFunction() UserInfoBatchFunctionType {
	return func(ctx context.Context, keys []UserInfoLoaderKey) []*UserInfoResultType {
		keysBySource := make(map[LoadUserInfoSource][]string)
		keyToIndexesMap := make(map[UserInfoLoaderKey][]int)

		for index, key := range keys {
			keysBySource[key.Source] = append(keysBySource[key.Source], key.PublicId)
			keyToIndexesMap[key] = append(keyToIndexesMap[key], index)
		}

		results := make([]*UserInfoResultType, len(keys))

		for source, publicIds := range keysBySource {
			var publicUserInfos []*gqlmodels.PublicUserInfo
			var exception *exceptions.Exception

			switch source {
			case LoadUserInfoSource_UserPublicId:
				// make sure we get the result in the same order
				// so the order of "publicUserInfos" is the same as the "publicIds"
				publicUserInfos, exception = d.userInfoService.GetPublicUserInfosByUserPublicIds(ctx, publicIds)
			default:
				exception = exceptions.UserInfo.InvalidSourceInBatchFunction()
			}

			if exception != nil {
				for _, publicId := range publicIds {
					key := UserInfoLoaderKey{PublicId: publicId, Source: source}
					if _, exists := keyToIndexesMap[key]; exists {
						for _, index := range keyToIndexesMap[key] {
							results[index] = &UserInfoResultType{Error: exception.Error}
						}
					}
				}
				continue
			}

			for index, publicUserInfo := range publicUserInfos {
				key := UserInfoLoaderKey{PublicId: publicIds[index], Source: source}
				if _, exists := keyToIndexesMap[key]; exists {
					for _, originalIndex := range keyToIndexesMap[key] {
						results[originalIndex] = &UserInfoResultType{Data: publicUserInfo}
					}
				}
			}
		}

		return results
	}
}

/* ============================== Load Functions ============================== */

func (d *UserInfoDataloader) LoadByUserPublicId(originalContext context.Context, publicId string) (*gqlmodels.PublicUserInfo, error) {
	future := d.loader.Load(
		originalContext,
		UserInfoLoaderKey{
			PublicId: publicId,
			Source:   LoadUserInfoSource_UserPublicId,
		},
	)

	publicUserInfo, err := future()
	if err != nil {
		return nil, err
	}

	return publicUserInfo, nil
}
