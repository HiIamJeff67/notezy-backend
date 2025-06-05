package enums

import (
	"database/sql/driver"
	"reflect"
)

/* ============================== Theme Definition ============================== */
type Theme string

const (
	Theme_Light  Theme = "Light"
	Theme_Dark   Theme = "Dark"
	Theme_System Theme = "System"
)

func (t Theme) Name() string {
	return reflect.TypeOf(t).Name()
}

func (t *Theme) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*t = Theme(string(v))
		return nil
	case string:
		*t = Theme(v)
		return nil
	}
	return scanError(value, t)
}

func (t Theme) Value() (driver.Value, error) {
	return string(t), nil
}

func (t *Theme) IsValidEnum() bool {
	for _, enum := range AllThemes {
		if *t == enum {
			return true
		}
	}
	return false
}

/* ========================= All Themes ========================= */
var AllThemes = []Theme{
	Theme_Light,
	Theme_Dark,
	Theme_System,
}
var AllThemeStrings = []string{
	string(Theme_Light),
	string(Theme_Dark),
	string(Theme_System),
}
