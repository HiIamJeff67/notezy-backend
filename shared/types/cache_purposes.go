package types

type ValidCachePurpose string

const (
	ValidCachePurpose_UserData    ValidCachePurpose = "UserData"
	ValidCachePurpose_RecentPages ValidCachePurpose = "RecentPages"
)

var _validCachePurposes = map[string]ValidCachePurpose{
	"UserData":    ValidCachePurpose_UserData,
	"RecentPages": ValidCachePurpose_RecentPages,
}

func (cp ValidCachePurpose) String() string {
	return string(cp)
}

func IsValidCachePurpose(cachePurpose string) bool {
	_, ok := _validCachePurposes[cachePurpose]
	return ok
}
func ConvertToValidCachePurpose(cachePurposeString string) (ValidCachePurpose, bool) {
	validCachePurpose, ok := _validCachePurposes[cachePurposeString]
	return validCachePurpose, ok
}
