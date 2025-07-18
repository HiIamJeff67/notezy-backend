package dataloaders

import (
	"gorm.io/gorm"
)

/* ============================== Interface & Instance ============================== */

type Dataloaders struct {
	UserDataLoader     UserDataloaderInterface
	UserInfoDataLoader UserInfoDataloaderInterface
	BadgeDataLoader    BadgeDataloaderInterface
}

func NewDataloaders(db *gorm.DB) Dataloaders {
	return Dataloaders{
		UserDataLoader:     NewUserDataloader(db),
		UserInfoDataLoader: NewUserInfoDataloader(db),
		BadgeDataLoader:    NewBadgeDataloader(db),
	}
}

/* ============================== General Methods for Other Dataloaders ============================== */

// func ClassifiableBatchFunction[LoaderKeyType any, BatchFunctionType any](ctx context.Context, keys []LoaderKeyType) BatchFunctionType {

// }
