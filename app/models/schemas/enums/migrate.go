package enums

// place the enums here to migrate
var MigratingEnums = map[string][]string{
	new(UserRole).Name():                AllUserRoleStrings,
	new(UserPlan).Name():                AllUserPlanStrings,
	new(UserStatus).Name():              AllUserStatusStrings,
	new(UserGender).Name():              AllUserGenderStrings,
	new(Country).Name():                 AllCountryStrings,
	new(CountryCode).Name():             AllCountryCodeStrings,
	new(Language).Name():                AllLanguageStrings,
	new(BadgeType).Name():               AllBadgeTypeStrings,
	new(AccessControlPermission).Name(): AllAccessControlPermissionStrings,
	new(MaterialType).Name():            AllMaterialTypeStrings,
	new(MaterialContentType).Name():     AllMaterialContentTypeStrings,
	new(SupportedBlockPackIcon).Name():  AllSupportedBlockPackIconStrings,
	new(BlockType).Name():               AllBlockTypeStrings,
}
