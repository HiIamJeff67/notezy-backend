package enums

import (
	"database/sql/driver"
	"fmt"
	"slices"
)

type Enum interface {
	Name() string
	Scan(value any) error
	Value() (driver.Value, error)
	IsValidEnum() bool
}

/* ==================== Temporary Function to Get the Scan Error ==================== */
func scanError(value any, e Enum) error {
	// A Helper Function to Get the Error
	return fmt.Errorf("failed to scan %T into %s", value, e.Name())
}

/* ========================= Validator for Validating Enums ========================= */
func IsValidEnumValues[EnumValue interface {
	UserRole |
		UserPlan |
		UserStatus |
		UserGender |
		Country |
		CountryCode |
		Theme |
		Language |
		BadgeType |
		string
}](value EnumValue, validateValues []EnumValue) bool {
	return slices.Contains(validateValues, value)
}

/* ========================= Map to Handling Migrating Enums ========================= */
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
