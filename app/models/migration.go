package models

import "notezy-backend/app/util"

var MigratingTables = []any{
	&User{},
	&UserInfo{},
	&UserAccount{},
	&UserSetting{},
	&UsersToBadges{},
	&Badge{},
}

var MigratingEnums = map[string][]string{
	util.GetTypeName[UserRole]():    AllUserRoleStrings,
	util.GetTypeName[UserPlan]():    AllUserPlanStrings,
	util.GetTypeName[UserStatus]():  AllUserStatusStrings,
	util.GetTypeName[UserGender]():  AllUserGenderStrings,
	util.GetTypeName[Country]():     AllCountryStrings,
	util.GetTypeName[CountryCode](): AllCountryCodeStrings,
	util.GetTypeName[Theme]():       AllThemeStrings,
	util.GetTypeName[Language]():    AllLanguageStrings,
	util.GetTypeName[BadgeType]():   AllBadgeTypeStrings,
}
