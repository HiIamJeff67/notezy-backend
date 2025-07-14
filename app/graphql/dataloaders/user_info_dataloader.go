package dataloaders

import (
	"context"
	"time"

	gophersdataloader "github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"

	gqlmodels "notezy-backend/app/graphql/models"
	services "notezy-backend/app/services"
)

type UserInfoDataloaderInterface interface {
	GetLoader() *gophersdataloader.Loader[string, *gqlmodels.PublicUserInfo]
	batchFunc() gophersdataloader.BatchFunc[string, *gqlmodels.PublicUserInfo]
}

type UserInfoDataloader struct {
	userInfoService services.UserInfoServiceInterface
	loader          *gophersdataloader.Loader[string, *gqlmodels.PublicUserInfo]
}

func NewUserInfoDataloader(db *gorm.DB) UserInfoDataloaderInterface {
	userInfoService := services.NewUserInfoService(db)
	dataloader := &UserInfoDataloader{
		userInfoService: userInfoService,
	}
	dataloader.loader = gophersdataloader.NewBatchedLoader(
		dataloader.batchFunc(),
		gophersdataloader.WithWait[string, *gqlmodels.PublicUserInfo](time.Microsecond),
	)

	return dataloader
}

func (d *UserInfoDataloader) GetLoader() *gophersdataloader.Loader[string, *gqlmodels.PublicUserInfo] {
	return d.loader
}

func (d *UserInfoDataloader) batchFunc() gophersdataloader.BatchFunc[string, *gqlmodels.PublicUserInfo] {
	return func(ctx context.Context, encodedSearchCursors []string) []*gophersdataloader.Result[*gqlmodels.PublicUserInfo] {
		publicUserInfos, exception := d.userInfoService.GetPublicUserInfosByEncodedSearchCursor(ctx, encodedSearchCursors)
		if exception != nil {
			results := make([]*gophersdataloader.Result[*gqlmodels.PublicUserInfo], len(encodedSearchCursors))
			for i := range results {
				results[i] = &gophersdataloader.Result[*gqlmodels.PublicUserInfo]{Error: exception.Error}
			}
			return results
		}

		results := make([]*gophersdataloader.Result[*gqlmodels.PublicUserInfo], len(encodedSearchCursors))
		for i, publicUserInfo := range publicUserInfos {
			results[i] = &gophersdataloader.Result[*gqlmodels.PublicUserInfo]{Data: publicUserInfo}
		}

		return results
	}
}
