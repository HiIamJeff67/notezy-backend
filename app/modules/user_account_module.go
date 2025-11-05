package modules

import (
	binders "notezy-backend/app/binders"
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
	userAccountRepository := repositories.NewUserAccountRepository()

	userAccountService := services.NewUserAccountService(
		models.NotezyDB,
		userAccountRepository,
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
