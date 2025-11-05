package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
)

type UserSettingModule struct {
	Binder     binders.UserSettingBinderInterface
	Controller controllers.UserSettingControllerInterface
}

func NewUserSettingModule() *UserSettingModule {
	userSettingRepository := repositories.NewUserSettingRepository()

	userSettingService := services.NewUserSettingService(
		models.NotezyDB,
		userSettingRepository,
	)

	userSettingBinder := binders.NewUserSettingBinder()

	userSettingController := controllers.NewUserSettingController(
		userSettingService,
	)

	return &UserSettingModule{
		Binder:     userSettingBinder,
		Controller: userSettingController,
	}
}
