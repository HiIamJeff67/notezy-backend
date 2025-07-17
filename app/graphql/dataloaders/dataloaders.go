package dataloaders

import (
	"gorm.io/gorm"
)

/* ============================== Interface & Instance ============================== */

type Dataloaders struct {
	UserInfoLoader UserInfoDataloaderInterface
	BadgeLoader    BadgeDataloaderInterface
}

func NewDataloaders(db *gorm.DB) Dataloaders {
	return Dataloaders{
		UserInfoLoader: NewUserInfoDataloader(db),
		BadgeLoader:    NewBadgeDataloader(db),
	}
}

/* ============================== General Methods for Other Dataloaders ============================== */

// func ClassifiableBatchFunction[LoaderKeyType any, BatchFunctionType any](ctx context.Context, keys []LoaderKeyType) BatchFunctionType {

// }
