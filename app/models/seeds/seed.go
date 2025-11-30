package seeds

import (
	_ "embed"

	planlimitationseeds "notezy-backend/app/models/seeds/plan_limitation_seeds"
)

var SeedingDefaultDataSQLs = []string{
	planlimitationseeds.PlanLimitationSeedingDefaultDataSQL_0000_UP,
}
