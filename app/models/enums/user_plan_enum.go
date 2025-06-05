package enums

import (
	"database/sql/driver"
	"reflect"
	"slices"
)

/* ============================== UserPlan Definition ============================== */
type UserPlan string

const (
	UserPlan_Enterprise UserPlan = "Enterprise"
	UserPlan_Ultimate   UserPlan = "Ultimate"
	UserPlan_Pro        UserPlan = "Pro"
	UserPlan_Free       UserPlan = "Free"
)

func (p UserPlan) Name() string {
	return reflect.TypeOf(p).Name()
}

func (p *UserPlan) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*p = UserPlan(string(v))
		return nil
	case string:
		*p = UserPlan(v)
		return nil
	}
	return scanError(value, p)
}

func (p UserPlan) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *UserPlan) IsValidEnum() bool {
	return slices.Contains(AllUserPlans, *p)
}

/* ========================= All UserPlans ========================= */
var AllUserPlans = []UserPlan{
	UserPlan_Enterprise,
	UserPlan_Ultimate,
	UserPlan_Pro,
	UserPlan_Free,
}
var AllUserPlanStrings = []string{
	string(UserPlan_Enterprise),
	string(UserPlan_Ultimate),
	string(UserPlan_Pro),
	string(UserPlan_Free),
}
