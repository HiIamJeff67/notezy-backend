package modules

import (
	binders "notezy-backend/app/binders"
	"notezy-backend/app/configs"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
)

type UserAccountModule struct {
	Binder     binders.UserAccountBinderInterface
	Controller controllers.UserAccountControllerInterface
}

func NewUserAccountModule() *UserAccountModule {
	userRepository := repositories.NewUserRepository()
	userAccountRepository := repositories.NewUserAccountRepository()
	oauthService := services.NewOAuthService(configs.OAuthGoogleConfig)

	userAccountService := services.NewUserAccountService(
		models.NotezyDB,
		userRepository,
		userAccountRepository,
		oauthService,
	)

	userAccountBinder := binders.NewUserAccountBinder()

	userAccountController := controllers.NewUserAccountController(
		userAccountService,
	)

	return &UserAccountModule{
		Binder:     userAccountBinder,
		Controller: userAccountController,
	}
}
