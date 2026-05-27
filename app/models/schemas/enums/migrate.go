package enums

// place the enums here to migrate
var MigratingEnums = map[string][]string{
	new(AccessControlPermission).Name():   AllAccessControlPermissionStrings,
	new(BadgeType).Name():                 AllBadgeTypeStrings,
	new(BillingIntervalUnit).Name():       AllBillingIntervalUnitStrings,
	new(BillingPlanName).Name():           AllBillingPlanNameStrings,
	new(BillingPlanStatus).Name():         AllBillingPlanStatusStrings,
	new(BlockType).Name():                 AllBlockTypeStrings,
	new(CountryCode).Name():               AllCountryCodeStrings,
	new(Country).Name():                   AllCountryStrings,
	new(ItemType).Name():                  AllItemTypeStrings,
	new(Language).Name():                  AllLanguageStrings,
	new(MaterialContentType).Name():       AllMaterialContentTypeStrings,
	new(RoutinePeriod).Name():             AllRoutinePeriodStrings,
	new(RoutineStatus).Name():             AllRoutineStatusStrings,
	new(RoutineTaskPurpose).Name():        AllRoutineTaskPurposeStrings,
	new(RoutineTaskStatus).Name():         AllRoutineTaskStatusStrings,
	new(SupportedIcon).Name():             AllSupportedIconStrings,
	new(SupportedCurrencyCode).Name():     AllSupportedCurrencyCodeStrings,
	new(UserGender).Name():                AllUserGenderStrings,
	new(UserPlan).Name():                  AllUserPlanStrings,
	new(UserRole).Name():                  AllUserRoleStrings,
	new(UserStatus).Name():                AllUserStatusStrings,
	new(UsersToBillingPlansStatus).Name(): AllUsersToBillingPlansStatusStrings,
}
