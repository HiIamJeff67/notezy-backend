package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

/* ============================== Definition ============================== */

type UserRole string

const (
	UserRole_Admin  UserRole = "Admin"
	UserRole_Normal UserRole = "Normal"
	UserRole_Guest  UserRole = "Guest"
)

/* ========================= All Instances ========================= */

// All the userRoles placing in the descending order
var AllUserRoles = []UserRole{
	UserRole_Admin,
	UserRole_Normal,
	UserRole_Guest,
}

// All the userRole strings placing in the descending order
var AllUserRoleStrings = []string{
	string(UserRole_Admin),
	string(UserRole_Normal),
	string(UserRole_Guest),
}

/* ============================== Methods ============================== */

func (r UserRole) Name() string {
	return reflect.TypeOf(r).Name()
}

// Scan() makes UserRole support automatically convert type from string in database to UserRole in codebase
func (r *UserRole) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*r = UserRole(string(v))
		return nil
	case string:
		*r = UserRole(v)
		return nil
	}
	return scanError(value, r)
}

// Value() makes UserRole support automatically convert from UserRole in codebase to string in database
func (r UserRole) Value() (driver.Value, error) {
	return string(r), nil
}

func (r UserRole) String() string {
	return string(r)
}

func (r *UserRole) IsValidEnum() bool {
	return slices.Contains(AllUserRoles, *r)
}

func ConvertStringToUserRole(enumString string) (*UserRole, error) {
	for _, userRole := range AllUserRoles {
		if string(userRole) == enumString {
			return &userRole, nil
		}
	}
	return nil, fmt.Errorf("invalid user role: %s", enumString)
}
