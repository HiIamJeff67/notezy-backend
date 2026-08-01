package dataloaders

import (
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type Dataloaders struct {
	UserDataLoader     UserDataloaderInterface
	UserInfoDataLoader UserInfoDataloaderInterface
	BadgeDataLoader    BadgeDataloaderInterface
}

func NewDataloaders(coreClient *coreadapters.CoreClient) Dataloaders {
	return Dataloaders{
		UserDataLoader:     NewUserDataloader(coreClient),
		UserInfoDataLoader: NewUserInfoDataloader(coreClient),
		BadgeDataLoader:    NewBadgeDataloader(coreClient),
	}
}
