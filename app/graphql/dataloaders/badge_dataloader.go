package dataloaders

import (
	"context"

	gophersdataloader "github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"

	"notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

/* ============================== Enum Keys & Types ============================== */

type LoadBadgeSource string

const (
	LoadBadgeSource_UserPublicId LoadBadgeSource = "LoadBadgeSourceUserPublicId"
)

type BadgeLoaderKey struct {
	ID     string          `json:"id"`
	Source LoadBadgeSource `json:"source"`
}

type BadgeLoaderType = gophersdataloader.Loader[BadgeLoaderKey, *gqlmodels.PublicBadge]
type BadgeBatchFunctionType = gophersdataloader.BatchFunc[BadgeLoaderKey, *gqlmodels.PublicBadge]
type BadgeResultType = gophersdataloader.Result[*gqlmodels.PublicBadge]

/* ============================== Interface & Instance ============================== */

type BadgeDataloaderInterface interface {
	GetLoader() *BadgeLoaderType
	batchFunction() BadgeBatchFunctionType

	// load functions
	LoadByUserPublicId(originalContext context.Context, id string) (*gqlmodels.PublicBadge, error)
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
		gophersdataloader.WithWait[BadgeLoaderKey, *gqlmodels.PublicBadge](constants.LoaderDelayOfBadge),
	)
	return dataloader
}

/* ============================== Dataloader Implementations ============================== */

func (d *BadgeDataloader) GetLoader() *BadgeLoaderType {
	return d.loader
}

// this batch function will fetch the PublicBadges using the publicIds of the "PublicUsers"
func (d *BadgeDataloader) batchFunction() BadgeBatchFunctionType {
	return func(ctx context.Context, keys []BadgeLoaderKey) []*BadgeResultType {
		keysBySource := make(map[LoadBadgeSource][]string)
		keyToIndexMap := make(map[BadgeLoaderKey]int)

		for index, key := range keys {
			keysBySource[key.Source] = append(keysBySource[key.Source], key.ID)
			keyToIndexMap[key] = index
		}

		results := make([]*BadgeResultType, len(keys))

		for source, ids := range keysBySource {
			var publicBadges []*gqlmodels.PublicBadge
			var exception *exceptions.Exception

			switch source {
			case LoadBadgeSource_UserPublicId:
				publicBadges, exception = d.badgeService.GetPublicBadgesByUserPublicIds(ctx, ids, true)
			default:
				exception = exceptions.Badge.NotFound().WithDetails(
					"Failed to fetch badge in a batch, source field is invalid",
				)
			}

			if exception != nil {
				for _, id := range ids {
					key := BadgeLoaderKey{ID: id, Source: source}
					if index, exists := keyToIndexMap[key]; exists {
						results[index] = &BadgeResultType{Error: exception.Error}
					}
				}
				continue
			}

			for index, publicBadge := range publicBadges {
				key := BadgeLoaderKey{ID: ids[index], Source: source}
				if originalIndex, exists := keyToIndexMap[key]; exists {
					results[originalIndex] = &BadgeResultType{Data: publicBadge}
				}
			}
		}

		return results
	}
}

/* ============================== Load Functions ============================== */

func (d *BadgeDataloader) LoadByUserPublicId(originalContext context.Context, id string) (*gqlmodels.PublicBadge, error) {
	future := d.loader.Load(
		originalContext,
		BadgeLoaderKey{
			ID:     id,
			Source: LoadBadgeSource_UserPublicId,
		},
	)

	publicBadge, err := future()
	if err != nil {
		return nil, err
	}

	return publicBadge, nil
}
