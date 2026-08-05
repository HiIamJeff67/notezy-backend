package dataloaders

import (
	"time"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

const (
	loaderDelayOfUser     = 50 * time.Microsecond
	loaderDelayOfUserInfo = 50 * time.Millisecond
	loaderDelayOfBadge    = 50 * time.Microsecond
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
