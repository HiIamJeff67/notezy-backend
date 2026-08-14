package userquotaconstraints

import _ "embed"

var (
	//go:embed drop_legacy_routine_task_cost_unit_count.sql
	DropLegacyRoutineTaskCostUnitCountSQL string
)
