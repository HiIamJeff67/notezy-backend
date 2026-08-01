package seeds

import (
	_ "embed"

	billingplanseeds "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/seeds/billing_plan_seeds"
	planlimitationseeds "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/seeds/plan_limitation_seeds"
)

var SeedingDefaultDataSQLs = []string{
	planlimitationseeds.PlanLimitationSeedingDefaultDataSQL_0000_UP,
	billingplanseeds.BillingPlanSeedingDefaultDataSQL_0000_UP,
}
