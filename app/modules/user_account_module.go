package modules

import (
	binders "github.com/HiIamJeff67/notezy-backend/app/binders"
	"github.com/HiIamJeff67/notezy-backend/app/configs"
	controllers "github.com/HiIamJeff67/notezy-backend/app/controllers"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
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
