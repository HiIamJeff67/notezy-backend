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
	ValidTableName_UsersToShelvesTable ValidTableName = "UsersToShelvesTable"
	ValidTableName_RootShelfTable      ValidTableName = "RootShelfTable"
	ValidTableName_SubShelfTable       ValidTableName = "SubShelfTable"
	ValidTableName_MaterialTable       ValidTableName = "MaterialTable"
	ValidTableName_BlockPackTable      ValidTableName = "BlockPackTable"
	ValidTableName_BlockGroupTable     ValidTableName = "BlockGroupTable"
	ValidTableName_BlockTable          ValidTableName = "BlockTable"
	ValidTableName_SyncBlockGroupTable ValidTableName = "SyncBlockGroupTable"
	ValidTableName_SyncBlockTableName  ValidTableName = "SyncBlockTable"
)

var _validTableNames = map[string]ValidTableName{
	"UserTable":           ValidTableName_UserTable,
	"UserAccountTable":    ValidTableName_UserAccountTable,
	"UserInfoTable":       ValidTableName_UserInfoTable,
	"UserSettingTable":    ValidTableName_UserSettingTable,
	"BadgeTable":          ValidTableName_BadgeTable,
	"UsersToBadgesTable":  ValidTableName_UsersToBadgesTable,
	"ThemeTable":          ValidTableName_ThemeTable,
	"UsersToShelvesTable": ValidTableName_UsersToShelvesTable,
	"RootShelfTable":      ValidTableName_RootShelfTable,
	"SubShelfTable":       ValidTableName_SubShelfTable,
	"MaterialTable":       ValidTableName_MaterialTable,
	"BlockPackTable":      ValidTableName_BlockPackTable,
	"BlockGroupTable":     ValidTableName_BlockGroupTable,
	"BlockTable":          ValidTableName_BlockTable,
	"SyncBlockGroupTable": ValidTableName_SyncBlockGroupTable,
	"SyncBlockTable":      ValidTableName_SyncBlockTableName,
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
