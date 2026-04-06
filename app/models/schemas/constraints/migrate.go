package constraints

import (
	blockgroupconstraints "notezy-backend/app/models/schemas/constraints/block_group_constraints"
	usesrtobillingplansconstraints "notezy-backend/app/models/schemas/constraints/usesr_to_billing_plans_constraints"
)

var MigratingConstraintSQLs = []string{
	blockgroupconstraints.BlockPackIdPrevBlockGroupIdIndexSQL,
	usesrtobillingplansconstraints.UserIdBillingPlanIdPartialStatusIndexSQL,
}
