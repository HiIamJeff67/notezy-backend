package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
)

type AuthModule struct {
	Binder     binders.AuthBinderInterface
	Controller controllers.AuthControllerInterface
}

func NewAuthModule() *AuthModule {
	userRepository := repositories.NewUserRepository()
	userInfoRepository := repositories.NewUserInfoRepository()
	userAccountRepository := repositories.NewUserAccountRepository()
	userSettingRepository := repositories.NewUserSettingRepository()

	authService := services.NewAuthService(
		models.NotezyDB,
		userRepository,
		userInfoRepository,
		userAccountRepository,
		userSettingRepository,
	)

	authBinder := binders.NewAuthBinder()

	authController := controllers.NewAuthController(
		authService,
	)

	return &AuthModule{
		Binder:     authBinder,
		Controller: authController,
	}
}
