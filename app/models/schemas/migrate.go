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
	&Block{},
	&Item{},

	&RoutinesToItems{},
	&UsersToStations{},
	&UsersToRoutineTags{},
	&Station{},
	&Routine{},
	&RoutinesToTasks{},
	&RoutineTask{},
	&RoutineTaskRecord{},
	&RoutinesToTags{},
	&RoutineTag{},

	&UsersToBillingPlans{},

	// private tables
	&PlanLimitation{},
	&BillingPlan{},
}
