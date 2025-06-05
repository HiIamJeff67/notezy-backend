package enums

import (
	"database/sql/driver"
	"reflect"
)

/* ============================== BadgeType Definition ============================== */
type BadgeType string

const (
	BadgeType_Diamond BadgeType = "Diamond"
	BadgeType_Golden  BadgeType = "Golden"
	BadgeType_Silver  BadgeType = "Silver"
	BadgeType_Bronze  BadgeType = "Bronze"
	BadgeType_Steel   BadgeType = "Steel"
)

func (bt BadgeType) Name() string {
	return reflect.TypeOf(bt).Name()
}

func (bt *BadgeType) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*bt = BadgeType(string(v))
		return nil
	case string:
		*bt = BadgeType(v)
		return nil
	}
	return scanError(value, bt)
}

func (bt BadgeType) Value() (driver.Value, error) {
	return string(bt), nil
}

func (bt *BadgeType) IsValidEnum() bool {
	for _, enum := range AllBadgeTypes {
		if *bt == enum {
			return true
		}
	}
	return false
}

/* ========================= All BadgeTypes ========================= */
var AllBadgeTypes = []BadgeType{
	BadgeType_Diamond,
	BadgeType_Golden,
	BadgeType_Silver,
	BadgeType_Bronze,
	BadgeType_Steel,
}
var AllBadgeTypeStrings = []string{
	string(BadgeType_Diamond),
	string(BadgeType_Golden),
	string(BadgeType_Silver),
	string(BadgeType_Bronze),
	string(BadgeType_Steel),
}
