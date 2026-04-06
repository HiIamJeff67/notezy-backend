package seeds

import (
	_ "embed"

	billingplanseeds "notezy-backend/app/models/seeds/billing_plan_seeds"
	planlimitationseeds "notezy-backend/app/models/seeds/plan_limitation_seeds"
)

var SeedingDefaultDataSQLs = []string{
	planlimitationseeds.PlanLimitationSeedingDefaultDataSQL_0000_UP,
	billingplanseeds.BillingPlanSeedingDefaultDataSQL_0000_UP,
}
