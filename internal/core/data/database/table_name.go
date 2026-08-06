package database

import platformdatabase "github.com/HiIamJeff67/notezy-backend/shared/platform/database"

const (
	TableName_UserTable        platformdatabase.TableName = "UserTable"
	TableName_UserAccountTable platformdatabase.TableName = "UserAccountTable"
	TableName_UserInfoTable    platformdatabase.TableName = "UserInfoTable"
	TableName_UserSettingTable platformdatabase.TableName = "UserSettingTable"

	TableName_BadgeTable         platformdatabase.TableName = "BadgeTable"
	TableName_UsersToBadgesTable platformdatabase.TableName = "UsersToBadgesTable"

	TableName_ThemeTable platformdatabase.TableName = "ThemeTable"

	TableName_UsersToShelvesTable       platformdatabase.TableName = "UsersToShelvesTable"
	TableName_RootShelfTable            platformdatabase.TableName = "RootShelfTable"
	TableName_SubShelfTable             platformdatabase.TableName = "SubShelfTable"
	TableName_MaterialTable             platformdatabase.TableName = "MaterialTable"
	TableName_BlockPackTable            platformdatabase.TableName = "BlockPackTable"
	TableName_BlockPackYjsDocumentTable platformdatabase.TableName = "BlockPackYjsDocumentTable"
	TableName_BlockPackYjsUpdateTable   platformdatabase.TableName = "BlockPackYjsUpdateTable"
	TableName_BlockTable                platformdatabase.TableName = "BlockTable"
	TableName_ItemTable                 platformdatabase.TableName = "ItemTable"

	TableName_RoutinesToItemsTable   platformdatabase.TableName = "RoutinesToItemsTable"
	TableName_UsersToStationsTable   platformdatabase.TableName = "UsersToStationsTable"
	TableName_StationTable           platformdatabase.TableName = "StationTable"
	TableName_RoutineTable           platformdatabase.TableName = "RoutineTable"
	TableName_RoutineDependencyTable platformdatabase.TableName = "RoutineDependencyTable"
	TableName_RoutineTaskTable       platformdatabase.TableName = "RoutineTaskTable"
	TableName_RoutineTaskRecordTable platformdatabase.TableName = "RoutineTaskRecordTable"
	TableName_RoutineTagTable        platformdatabase.TableName = "RoutineTagTable"
	TableName_RoutinesToTagsTable    platformdatabase.TableName = "RoutinesToTagsTable"

	TableName_UsersToBillingPlansTable platformdatabase.TableName = "UsersToBillingPlansTable"

	TableName_PlanLimitationTable platformdatabase.TableName = "PlanLimitationTable"
	TableName_BillingPlanTable    platformdatabase.TableName = "BillingPlanTable"
)

var _validTableNames = map[string]platformdatabase.TableName{
	"UserTable":        TableName_UserTable,
	"UserAccountTable": TableName_UserAccountTable,
	"UserInfoTable":    TableName_UserInfoTable,
	"UserSettingTable": TableName_UserSettingTable,

	"BadgeTable":         TableName_BadgeTable,
	"UsersToBadgesTable": TableName_UsersToBadgesTable,

	"ThemeTable": TableName_ThemeTable,

	"UsersToShelvesTable":       TableName_UsersToShelvesTable,
	"RootShelfTable":            TableName_RootShelfTable,
	"SubShelfTable":             TableName_SubShelfTable,
	"MaterialTable":             TableName_MaterialTable,
	"BlockPackTable":            TableName_BlockPackTable,
	"BlockPackYjsDocumentTable": TableName_BlockPackYjsDocumentTable,
	"BlockPackYjsUpdateTable":   TableName_BlockPackYjsUpdateTable,
	"BlockTable":                TableName_BlockTable,
	"ItemTable":                 TableName_ItemTable,

	"RoutinesToItemsTable":   TableName_RoutinesToItemsTable,
	"UsersToStationsTable":   TableName_UsersToStationsTable,
	"StationTable":           TableName_StationTable,
	"RoutineTable":           TableName_RoutineTable,
	"RoutineDependencyTable": TableName_RoutineDependencyTable,
	"RoutineTaskTable":       TableName_RoutineTaskTable,
	"RoutineTaskRecordTable": TableName_RoutineTaskRecordTable,
	"RoutineTagTable":        TableName_RoutineTagTable,
	"RoutinesToTagsTable":    TableName_RoutinesToTagsTable,

	"UsersToBillingPlansTable": TableName_UsersToBillingPlansTable,

	"PlanLimitationTable": TableName_PlanLimitationTable,
	"BillingPlanTable":    TableName_BillingPlanTable,
}

func ConvertToTableName(tableName string) (platformdatabase.TableName, bool) {
	validTableName, ok := _validTableNames[tableName]
	return validTableName, ok
}
