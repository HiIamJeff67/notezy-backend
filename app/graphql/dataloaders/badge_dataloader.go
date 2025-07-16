package dataloaders

import (
	"context"

	gophersdataloader "github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"

	gqlmodels "notezy-backend/app/graphql/models"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

type BadgeLoaderType = gophersdataloader.Loader[string, *gqlmodels.PublicBadge]
type BadgeBatchFunctionType = gophersdataloader.BatchFunc[string, *gqlmodels.PublicBadge]
type BadgeResultType = gophersdataloader.Result[*gqlmodels.PublicBadge]

type BadgeDataloaderInterface interface {
	GetLoader() *BadgeLoaderType
	batchFunction() BadgeBatchFunctionType
}

type BadgeDataloader struct {
	badgeService services.BadgeServiceInterface
	loader       *BadgeLoaderType
}

func NewBadgeDataloader(db *gorm.DB) BadgeDataloaderInterface {
	dataloader := &BadgeDataloader{
		badgeService: services.NewBadgeService(db),
	}
	dataloader.loader = gophersdataloader.NewBatchedLoader(
		dataloader.batchFunction(),
		gophersdataloader.WithWait[string, *gqlmodels.PublicBadge](constants.LoaderDelayOfBadge),
	)
	return dataloader
}

func (d *BadgeDataloader) GetLoader() *BadgeLoaderType {
	return d.loader
}

// this batch function will fetch the PublicBadges using the publicIds of the "PublicUsers"
func (d *BadgeDataloader) batchFunction() BadgeBatchFunctionType {
	return func(ctx context.Context, publicIds []string) []*BadgeResultType {
		publicBadges, exception := d.badgeService.GetPublicBadgesByPublicUserIds(ctx, publicIds)
		if exception != nil {
			results := make([]*BadgeResultType, len(publicIds))
			for i := range results {
				results[i] = &gophersdataloader.Result[*gqlmodels.PublicBadge]{Error: exception.Error}
			}
			return results
		}

		results := make([]*BadgeResultType, len(publicIds))
		for i, publicBadge := range publicBadges {
			results[i] = &gophersdataloader.Result[*gqlmodels.PublicBadge]{Data: publicBadge}
		}

		return results
	}
}
