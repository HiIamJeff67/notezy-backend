package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
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
