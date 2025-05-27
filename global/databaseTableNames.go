package global

type ValidTableName string

const (
	ValidTableName_UserTable ValidTableName = "UserTable"
	ValidTableName_UserAccountTable ValidTableName = "UserAccountTable"
	ValidTableName_UserInfoTable ValidTableName = "UserInfoTable"
	ValidTableName_UserSettingTable ValidTableName = "UserSettingTable"
	ValidTableName_BadgeTable ValidTableName = "BadgeTable"
	ValidTableName_UsersToBadgesTable ValidTableName = "UsersToBadgesTable"
)
var _validTableNames = map[string]ValidTableName{
	"UserTable": ValidTableName_UserTable, 
	"UserAccountTable": ValidTableName_UserAccountTable, 
	"UserInfoTable": ValidTableName_UserInfoTable, 
	"UserSettingTable": ValidTableName_UserSettingTable,
	"BadgeTable": ValidTableName_BadgeTable, 
	"UsersToBadgesTable": ValidTableName_UsersToBadgesTable, 
}
func IsValidTableName(tableName string) bool {
	_, ok := _validTableNames[tableName]
	return ok
}
func ConvertToValidTableName(tableName string) (ValidTableName, bool) {
	validTableName, ok := _validTableNames[tableName]
	return validTableName, ok
}