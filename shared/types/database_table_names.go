package types

type TableName string

const (
	// public tables(accessable and mutatable by the client user and admin)
	TableName_UserTable           TableName = "UserTable"
	TableName_UserAccountTable    TableName = "UserAccountTable"
	TableName_UserInfoTable       TableName = "UserInfoTable"
	TableName_UserSettingTable    TableName = "UserSettingTable"
	TableName_BadgeTable          TableName = "BadgeTable"
	TableName_UsersToBadgesTable  TableName = "UsersToBadgesTable"
	TableName_ThemeTable          TableName = "ThemeTable"
	TableName_UsersToShelvesTable TableName = "UsersToShelvesTable"
	TableName_RootShelfTable      TableName = "RootShelfTable"
	TableName_SubShelfTable       TableName = "SubShelfTable"
	TableName_MaterialTable       TableName = "MaterialTable"
	TableName_BlockPackTable      TableName = "BlockPackTable"
	TableName_BlockGroupTable     TableName = "BlockGroupTable"
	TableName_BlockTable          TableName = "BlockTable"
	TableName_SyncBlockGroupTable TableName = "SyncBlockGroupTable"
	TableName_SyncBlockTable      TableName = "SyncBlockTable"

	// private tables(accessable by the client user and admin, but only mutatable by the admin)
	TableName_PlanLimitationTable TableName = "PlanLimitationTable"
)

var _validTableNames = map[string]TableName{
	// public tables
	"UserTable":           TableName_UserTable,
	"UserAccountTable":    TableName_UserAccountTable,
	"UserInfoTable":       TableName_UserInfoTable,
	"UserSettingTable":    TableName_UserSettingTable,
	"BadgeTable":          TableName_BadgeTable,
	"UsersToBadgesTable":  TableName_UsersToBadgesTable,
	"ThemeTable":          TableName_ThemeTable,
	"UsersToShelvesTable": TableName_UsersToShelvesTable,
	"RootShelfTable":      TableName_RootShelfTable,
	"SubShelfTable":       TableName_SubShelfTable,
	"MaterialTable":       TableName_MaterialTable,
	"BlockPackTable":      TableName_BlockPackTable,
	"BlockGroupTable":     TableName_BlockGroupTable,
	"BlockTable":          TableName_BlockTable,
	"SyncBlockGroupTable": TableName_SyncBlockGroupTable,
	"SyncBlockTable":      TableName_SyncBlockTable,

	// private tables
	"PlanLimitationTable": TableName_PlanLimitationTable,
}

func (tn TableName) String() string {
	return string(tn)
}

func IsTableName(tableName string) bool {
	_, ok := _validTableNames[tableName]
	return ok
}
func ConvertToTableName(tableName string) (TableName, bool) {
	validTableName, ok := _validTableNames[tableName]
	return validTableName, ok
}
