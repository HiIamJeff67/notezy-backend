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
	&UsersToBadges{},
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

	// private tables
	&PlanLimitation{},
}
