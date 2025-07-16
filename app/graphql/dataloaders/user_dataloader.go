package dataloaders

import (
	gqlmodels "notezy-backend/app/graphql/models"
	"notezy-backend/app/services"

	gophersdataloader "github.com/graph-gophers/dataloader/v7"
)

type UserLoaderType = gophersdataloader.Loader[string, *gqlmodels.PublicUser]
type UserBatchFunctionType = gophersdataloader.BatchFunc[string, *gqlmodels.PublicUser]
type UserResultType = gophersdataloader.Result[*gqlmodels.PublicUser]

type UserDataloaderInterface interface{}

type UserDataloader struct {
	userService services.UserServiceInterface
	loader      *UserLoaderType
}

// func NewUserDataloader(db *gorm.DB) UserDataloaderInterface {
// 	dataloader := &UserDataloader{
// 		userService: services.NewUserService(db),
// 	}
// 	datalodader.loader = gophersdataloader.NewBatchedLoader(
// 		dataloader.batchFun
// 	)
// }

func (d *UserDataloader) GetLoader() *UserLoaderType {
	return d.loader
}

// func (d *UserDataloader) batchFunction() UserBatchFunctionType {
// 	return func(ctx context.Context, publicIds []string) []*UserResultType {
// 		publicUsers, exception := d.userService.GetPublicUserByPublicId(ctx, publicIds)
// 		if exception != nil {
// 			results := make([]*UserResultType, len(publicIds))
// 			for i := range results {
// 				results[i] = &UserResultType{Error: exception.Error}
// 			}
// 			return results
// 		}

// 		results := make([]*UserResultType, len(publicIds))
// 		for i, publicUser := range publicUsers {
// 			results[i] = &UserResultType{Data: publicUser}
// 		}

// 		return results
// 	}
// }
