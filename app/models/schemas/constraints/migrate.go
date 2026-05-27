package constraints

import (
	blockconstraints "notezy-backend/app/models/schemas/constraints/block_constraints"
	routineconstraints "notezy-backend/app/models/schemas/constraints/routine_constraints"
	userstobillingplansconstraints "notezy-backend/app/models/schemas/constraints/users_to_billing_plans_constraints"
)

var MigratingConstraintSQLs = []string{
	userstobillingplansconstraints.UserIdBillingPlanIdPartialStatusIndexSQL,
	blockconstraints.BlockTreeRootIndexSQL,
	routineconstraints.RoutineScheduledTimeMinutePrecisionCheckSQL,
	routineconstraints.RoutineScheduledTimeInPeriodCheckSQL,
}
