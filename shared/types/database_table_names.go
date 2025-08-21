package types

type ValidTableName string

const (
	ValidTableName_UserTable           ValidTableName = "UserTable"
	ValidTableName_UserAccountTable    ValidTableName = "UserAccountTable"
	ValidTableName_UserInfoTable       ValidTableName = "UserInfoTable"
	ValidTableName_UserSettingTable    ValidTableName = "UserSettingTable"
	ValidTableName_BadgeTable          ValidTableName = "BadgeTable"
	ValidTableName_UsersToBadgesTable  ValidTableName = "UsersToBadgesTable"
	ValidTableName_ThemeTable          ValidTableName = "ThemeTable"
	ValidTableName_ShelfTable          ValidTableName = "ShelfTable"
	ValidTableName_UsersToShelvesTable ValidTableName = "UsersToShelvesTable"
	ValidTableName_MaterialTable       ValidTableName = "MaterialTable"
)

var _validTableNames = map[string]ValidTableName{
	"UserTable":           ValidTableName_UserTable,
	"UserAccountTable":    ValidTableName_UserAccountTable,
	"UserInfoTable":       ValidTableName_UserInfoTable,
	"UserSettingTable":    ValidTableName_UserSettingTable,
	"BadgeTable":          ValidTableName_BadgeTable,
	"UsersToBadgesTable":  ValidTableName_UsersToBadgesTable,
	"ThemeTable":          ValidTableName_ThemeTable,
	"ShelfTable":          ValidTableName_ShelfTable,
	"UsersToShelvesTable": ValidTableName_UsersToShelvesTable,
	"MaterialTable":       ValidTableName_MaterialTable,
}

func (tn ValidTableName) String() string {
	return string(tn)
}

func IsValidTableName(tableName string) bool {
	_, ok := _validTableNames[tableName]
	return ok
}
func ConvertToValidTableName(tableName string) (ValidTableName, bool) {
	validTableName, ok := _validTableNames[tableName]
	return validTableName, ok
}
