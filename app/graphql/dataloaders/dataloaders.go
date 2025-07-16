package dataloaders

import (
	"gorm.io/gorm"
)

type Dataloaders struct {
	UserInfoLoader *UserInfoLoaderType
	BadgeLoader    *BadgeLoaderType
}

func NewDataloaders(db *gorm.DB) Dataloaders {
	return Dataloaders{
		UserInfoLoader: NewUserInfoDataloader(db).GetLoader(),
		BadgeLoader:    NewBadgeDataloader(db).GetLoader(),
	}
}
