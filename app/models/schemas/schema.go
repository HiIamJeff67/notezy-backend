package schemas

// place the tables here to migrate
var MigratingTables = []any{
	&User{},
	&UserInfo{},
	&UserAccount{},
	&UserSetting{},
	&UsersToBadges{},
	&Badge{},
	&UsersToBadges{},
	&Theme{},
}
