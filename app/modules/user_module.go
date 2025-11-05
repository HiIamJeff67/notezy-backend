package modules

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	models "notezy-backend/app/models"
	repositories "notezy-backend/app/models/repositories"
	services "notezy-backend/app/services"
)

type UserModule struct {
	Binder     binders.UserBinderInterface
	Controller controllers.UserControllerInterface
}

func NewUserModule() *UserModule {
	userRepository := repositories.NewUserRepository()

	userService := services.NewUserService(
		models.NotezyDB,
		userRepository,
	)

	userBinder := binders.NewUserBinder()

	userController := controllers.NewUserController(
		userService,
	)

	return &UserModule{
		Binder:     userBinder,
		Controller: userController,
	}
}
