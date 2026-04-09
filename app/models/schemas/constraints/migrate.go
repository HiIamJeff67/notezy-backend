package constraints

import (
	usesrtobillingplansconstraints "notezy-backend/app/models/schemas/constraints/usesr_to_billing_plans_constraints"
)

var MigratingConstraintSQLs = []string{
	usesrtobillingplansconstraints.UserIdBillingPlanIdPartialStatusIndexSQL,
}
