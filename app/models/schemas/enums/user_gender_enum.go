package enums

import (
	"database/sql/driver"
	"reflect"
	"slices"
)

/* ============================== UserGener Definition ============================== */
type UserGender string

const (
	UserGender_Male           UserGender = "Male"
	UserGender_Female         UserGender = "Female"
	UserGender_PreferNotToSay UserGender = "PreferNotToSay"
)

func (g UserGender) Name() string {
	return reflect.TypeOf(g).Name()
}

func (g *UserGender) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*g = UserGender(string(v))
		return nil
	case string:
		*g = UserGender(v)
		return nil
	}
	return scanError(value, g)
}

func (g UserGender) Value() (driver.Value, error) {
	return string(g), nil
}

func (g UserGender) String() string {
	return string(g)
}

func (g *UserGender) IsValidEnum() bool {
	return slices.Contains(AllUserGenders, *g)
}

/* ========================= All UserGenders ========================= */
var AllUserGenders = []UserGender{
	UserGender_Male,
	UserGender_Female,
	UserGender_PreferNotToSay,
}
var AllUserGenderStrings = []string{
	string(UserGender_Male),
	string(UserGender_Female),
	string(UserGender_PreferNotToSay),
}
