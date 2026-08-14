package dataloaders

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	gophersdataloader "github.com/graph-gophers/dataloader/v7"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-infos"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/graphql/models"

	gatewaycontexts "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/contexts"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

type LoadUserInfoSource string

const (
	LoadUserInfoSourceUserPublicId LoadUserInfoSource = "LoadUserInfoSourceUserPublicId"
)

type UserInfoLoaderKey struct {
	PublicId uuid.UUID          `json:"publicId"`
	Source   LoadUserInfoSource `json:"source"`
}

type UserInfoLoaderType = gophersdataloader.Loader[UserInfoLoaderKey, *gqlmodels.PublicUserInfo]
type UserInfoBatchFunctionType = gophersdataloader.BatchFunc[UserInfoLoaderKey, *gqlmodels.PublicUserInfo]
type UserInfoResultType = gophersdataloader.Result[*gqlmodels.PublicUserInfo]

type UserInfoDataloaderInterface interface {
	GetLoader() *UserInfoLoaderType
	LoadByUserPublicId(ctx context.Context, publicId uuid.UUID) (*gqlmodels.PublicUserInfo, error)
}

type UserInfoDataloader struct {
	coreClient *coreadapters.CoreAdapter
	loader     *UserInfoLoaderType
}

func NewUserInfoDataloader(coreClient *coreadapters.CoreAdapter) UserInfoDataloaderInterface {
	dataloader := &UserInfoDataloader{
		coreClient: coreClient,
	}
	dataloader.loader = gophersdataloader.NewBatchedLoader(
		dataloader.batchFunction(),
		gophersdataloader.WithWait[UserInfoLoaderKey, *gqlmodels.PublicUserInfo](loaderDelayOfUserInfo),
	)

	return dataloader
}

/* ============================== Dataloader Methods ============================== */

func (d *UserInfoDataloader) GetLoader() *UserInfoLoaderType {
	return d.loader
}

func (d *UserInfoDataloader) batchFunction() UserInfoBatchFunctionType {
	return func(ctx context.Context, keys []UserInfoLoaderKey) []*UserInfoResultType {
		results := make([]*UserInfoResultType, len(keys))
		publicIds := make([]uuid.UUID, 0, len(keys))
		indexesByPublicId := make(map[uuid.UUID][]int, len(keys))

		for index, key := range keys {
			if key.Source != LoadUserInfoSourceUserPublicId {
				exception := exceptions.New(
					"InvalidSource",
					"GraphQL",
					"LoadUserInfos",
					"UserInfo dataloader source is invalid",
					http.StatusInternalServerError,
					true,
				)
				results[index] = &UserInfoResultType{
					Error: exception.Origin(),
				}
				continue
			}

			publicIds = append(publicIds, key.PublicId)
			indexesByPublicId[key.PublicId] = append(indexesByPublicId[key.PublicId], index)
		}
		if len(publicIds) == 0 {
			return results
		}

		ginContext, exception := gatewaycontexts.GetAndConvertContextToGinContext(ctx)
		if exception != nil {
			for _, indexes := range indexesByPublicId {
				for _, index := range indexes {
					results[index] = &UserInfoResultType{
						Error: exception.Origin(),
					}
				}
			}
			return results
		}

		requestDto := apicontract.LoadUserInfosRequestDto(publicIds)
		response, exception := coreadapters.CallSecurly[
			apicontract.LoadUserInfosRequestDto,
			apicontract.LoadUserInfosResponseDto,
		](
			ginContext,
			d.coreClient,
			&requestDto,
			apicontract.LoadUserInfosOperation,
			"/core/v1/user-infos/graphql/load",
		)
		if exception != nil {
			for _, indexes := range indexesByPublicId {
				for _, index := range indexes {
					results[index] = &UserInfoResultType{
						Error: exception.Origin(),
					}
				}
			}
			return results
		}

		for index, model := range response.Data {
			if index >= len(publicIds) {
				break
			}
			for _, resultIndex := range indexesByPublicId[publicIds[index]] {
				results[resultIndex] = &UserInfoResultType{
					Data: model,
				}
			}
		}

		return results
	}
}

func (d *UserInfoDataloader) LoadByUserPublicId(
	ctx context.Context,
	publicId uuid.UUID,
) (*gqlmodels.PublicUserInfo, error) {
	future := d.loader.Load(ctx, UserInfoLoaderKey{
		PublicId: publicId,
		Source:   LoadUserInfoSourceUserPublicId,
	})

	return future()
}
