package schemas

// place the tables here to migrate
var MigratingTables = []any{
	// public tables
	&User{},
	&UserInfo{},
	&UserAccount{},
	&UserQuota{},
	&UserSetting{},
	&APIKey{},

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
	&InboxEvent{},
	&OutboxEvent{},

	&UsersToBillingPlans{},

	// private tables
	&PlanLimitation{},
	&BillingPlan{},
}
