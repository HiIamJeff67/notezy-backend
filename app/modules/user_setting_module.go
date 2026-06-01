package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
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
