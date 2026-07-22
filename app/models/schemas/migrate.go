package schemas

// place the tables here to migrate
var MigratingTables = []any{
	// public tables
	&User{},
	&UserInfo{},
	&UserAccount{},
	&UserSetting{},

	&UsersToBadges{},
	&Badge{},

	&Theme{},

	&UsersToShelves{},
	&RootShelf{},
	&SubShelf{},
	&Material{},
	&BlockPack{},
	&BlockPackYjsDocument{},
	&BlockPackYjsUpdate{},
	&Block{},
	&Item{},

	&Station{},
	&Routine{},
	&RoutineTag{},
	&UsersToStations{},
	&RoutinesToItems{},
	&RoutinesToTags{},
	&RoutineTask{},
	&RoutineTaskRecord{},

	&UsersToBillingPlans{},

	// private tables
	&PlanLimitation{},
	&BillingPlan{},
}
