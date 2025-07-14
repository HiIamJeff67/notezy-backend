package dataloaders

import (
	gophersdataloader "github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"

	gqlmodels "notezy-backend/app/graphql/models"
)

type Dataloaders struct {
	UserInfoDataloader *gophersdataloader.Loader[string, *gqlmodels.PublicUserInfo]
}

func NewDataloaders(db *gorm.DB) Dataloaders {
	return Dataloaders{
		UserInfoDataloader: NewUserInfoDataloader(db).GetLoader(),
	}
}
