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
	String() string
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
		Language |
		BadgeType |
		AccessControlPermission |
		MaterialType |
		MaterialContentType |
		BlockType |
		SupportedBlockPackIcon |
		string
}](value EnumValue, validateValues []EnumValue) bool {
	return slices.Contains(validateValues, value)
}
