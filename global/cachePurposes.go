package global

type ValidCachePurpose string

const (
	ValidCachePurpose_UserData ValidCachePurpose = "UserData"
	ValidCachePurpose_RecentPages ValidCachePurpose = "RecentPages"
)
var _validCachePurposes = map[string]ValidCachePurpose{
	"UserData": ValidCachePurpose_UserData, 
	"RecentPages": ValidCachePurpose_RecentPages, 
}
func IsValidCachePurpose(cachePurpose string) bool {
	_, ok := _validCachePurposes[cachePurpose]
	return ok
}
func ConvertToValidCachePurpose(cachePurpose string) (ValidCachePurpose, bool) {
	validCachePurpose, ok := _validCachePurposes[cachePurpose]
	return validCachePurpose, ok
}