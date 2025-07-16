package dataloaders

import (
	"context"

	gophersdataloader "github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"

	gqlmodels "notezy-backend/app/graphql/models"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

type UserInfoLoaderType = gophersdataloader.Loader[string, *gqlmodels.PublicUserInfo]
type UserInfoBatchFunctionType = gophersdataloader.BatchFunc[string, *gqlmodels.PublicUserInfo]
type UserInfoResultType = gophersdataloader.Result[*gqlmodels.PublicUserInfo]

type UserInfoDataloaderInterface interface {
	GetLoader() *UserInfoLoaderType
	batchFunction() UserInfoBatchFunctionType
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
		gophersdataloader.WithWait[string, *gqlmodels.PublicUserInfo](constants.LoaderDelayOfUserInfo),
	)

	return dataloader
}

func (d *UserInfoDataloader) GetLoader() *UserInfoLoaderType {
	return d.loader
}

// this batch function will fetch the PublicUserInfos using the publicIds of the "PublicUsers"
func (d *UserInfoDataloader) batchFunction() UserInfoBatchFunctionType {
	return func(ctx context.Context, publicIds []string) []*UserInfoResultType {
		publicUserInfos, exception := d.userInfoService.GetPublicUserInfosByPublicIds(ctx, publicIds)
		if exception != nil {
			results := make([]*UserInfoResultType, len(publicIds))
			for i := range results {
				results[i] = &UserInfoResultType{Error: exception.Error}
			}
			return results
		}

		results := make([]*UserInfoResultType, len(publicIds))
		for i, publicUserInfo := range publicUserInfos {
			results[i] = &UserInfoResultType{Data: publicUserInfo}
		}

		return results
	}
}
