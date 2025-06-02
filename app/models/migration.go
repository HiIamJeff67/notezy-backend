package models

// place the tables here to migrate
var MigratingTables = []any{
	&User{},
	&UserInfo{},
	&UserAccount{},
	&UserSetting{},
	&UsersToBadges{},
	&Badge{},
}

// place the enums here to migrate
var MigratingEnums = map[string][]string{
	new(UserRole).Name():    AllUserRoleStrings,
	new(UserPlan).Name():    AllUserPlanStrings,
	new(UserStatus).Name():  AllUserStatusStrings,
	new(UserGender).Name():  AllUserGenderStrings,
	new(Country).Name():     AllCountryStrings,
	new(CountryCode).Name(): AllCountryCodeStrings,
	new(Theme).Name():       AllThemeStrings,
	new(Language).Name():    AllLanguageStrings,
	new(BadgeType).Name():   AllBadgeTypeStrings,
}
