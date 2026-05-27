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
	&BlockGroup{},
	&Block{},
	&SyncBlockGroup{},
	&SyncBlock{},
	&Item{},

	&UsersToStations{},
	&Station{},
	&Routine{},
	&RoutinesToTasks{},
	&RoutineTask{},
	&RoutinesToTags{},
	&RoutineTag{},
	&RoutinesToItems{},

	&UsersToBillingPlans{},

	// private tables
	&PlanLimitation{},
	&BillingPlan{},
}
