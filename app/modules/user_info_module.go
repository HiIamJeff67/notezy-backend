package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type UserInfoModule struct {
	Binder     binders.UserInfoBinderInterface
	Controller controllers.UserInfoControllerInterface
}

func NewUserInfoModule() *UserInfoModule {
	userInfoRepository := repositories.NewUserInfoRepository()

	userInfoService := services.NewUserInfoService(
		models.NotezyDB,
		userInfoRepository,
	)

	userInfoBinder := binders.NewUserInfoBinder()

	userInfoController := controllers.NewUserInfoController(
		userInfoService,
	)

	return &UserInfoModule{
		Binder:     userInfoBinder,
		Controller: userInfoController,
	}
}
