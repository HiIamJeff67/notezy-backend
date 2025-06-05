package global

type ValidCookieName string

const (
	ValidCookieName_AccessToken  ValidCookieName = "accessToken"
	ValidCookieName_RefreshToken ValidCookieName = "refreshToken"
)

var _validCookieNames = map[string]ValidCookieName{
	"accessToken":  ValidCookieName_AccessToken,
	"refreshToken": ValidCookieName_RefreshToken,
}

func (cn ValidCookieName) String() string {
	return string(cn)
}

func IsValidCookieName(validCookieName string) bool {
	_, ok := _validCookieNames[validCookieName]
	return ok
}
func ConvertToValidCookieName(cachePurposeString string) (ValidCookieName, bool) {
	validCookieName, ok := _validCookieNames[cachePurposeString]
	return validCookieName, ok
}
