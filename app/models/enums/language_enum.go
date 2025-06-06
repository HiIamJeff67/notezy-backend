package enums

import (
	"database/sql/driver"
	"reflect"
	"slices"
)

/* ============================== Language Definition ============================== */
type Language string

const (
	Language_English            Language = "English"
	Language_TraditionalChinese Language = "TraditionalChinese"
	Language_SimpleChinese      Language = "SimpleChinese"
	Language_Japanese           Language = "Japanese"
)

func (l Language) Name() string {
	return reflect.TypeOf(l).Name()
}

func (l *Language) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*l = Language(string(v))
		return nil
	case string:
		*l = Language(v)
		return nil
	}
	return scanError(value, l)
}

func (l Language) Value() (driver.Value, error) {
	return string(l), nil
}

func (l Language) String() string {
	return string(l)
}

func (l *Language) IsValidEnum() bool {
	return slices.Contains(AllLanguages, *l)
}

/* ========================= All Languages ========================= */
var AllLanguages = []Language{
	Language_English,
	Language_TraditionalChinese,
	Language_SimpleChinese,
	Language_Japanese,
}
var AllLanguageStrings = []string{
	string(Language_English),
	string(Language_TraditionalChinese),
	string(Language_SimpleChinese),
	string(Language_Japanese),
}
