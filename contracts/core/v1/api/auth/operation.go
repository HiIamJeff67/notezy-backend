package apicontract

const (
	RegisterOperation          = "auth.register"
	RegisterViaGoogleOperation = "auth.register-via-google"
	LoginOperation             = "auth.login"
	LoginViaGoogleOperation    = "auth.login-via-google"
	LogoutOperation            = "auth.logout"
	SendAuthCodeOperation      = "auth.send-auth-code"
	ValidateEmailOperation     = "auth.validate-email"
	ResetEmailOperation        = "auth.reset-email"
	ForgetPasswordOperation    = "auth.forget-password"
	ResetMeOperation           = "auth.reset-me"
	DeleteMeOperation          = "auth.delete-me"
)
