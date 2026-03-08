package types

type ValidCachePurpose string

const (
	ValidCachePurpose_UserData   ValidCachePurpose = "UserData"
	ValidCachePurpose_RateLimite ValidCachePurpose = "RateLimit"
)

var AllValidCachePurposes = []ValidCachePurpose{
	ValidCachePurpose_UserData,
	ValidCachePurpose_RateLimite,
}

var _validCachePurposes = map[string]ValidCachePurpose{
	"UserData":  ValidCachePurpose_UserData,
	"RateLimit": ValidCachePurpose_RateLimite,
}

func (cp ValidCachePurpose) String() string {
	return string(cp)
}

func IsValidCachePurpose(cachePurposeString string) bool {
	_, ok := _validCachePurposes[cachePurposeString]
	return ok
}

func ConvertToValidCachePurpose(cachePurposeString string) (ValidCachePurpose, bool) {
	validCachePurpose, ok := _validCachePurposes[cachePurposeString]
	return validCachePurpose, ok
}
