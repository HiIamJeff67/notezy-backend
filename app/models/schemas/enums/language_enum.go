package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

/* ============================== Definition ============================== */

type Language string

const (
	Language_English            Language = "English"
	Language_TraditionalChinese Language = "TraditionalChinese"
	Language_SimpleChinese      Language = "SimpleChinese"
	Language_Japanese           Language = "Japanese"
	Language_Korean             Language = "Korean"
)

/* ========================= All Instances ========================= */

var AllLanguages = []Language{
	Language_English,
	Language_TraditionalChinese,
	Language_SimpleChinese,
	Language_Japanese,
	Language_Korean,
}
var AllLanguageStrings = []string{
	string(Language_English),
	string(Language_TraditionalChinese),
	string(Language_SimpleChinese),
	string(Language_Japanese),
	string(Language_Korean),
}

/* ============================== Methods ============================== */

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

func ConvertStringToLanguage(enumString string) (*Language, error) {
	for _, language := range AllLanguages {
		if string(language) == enumString {
			return &language, nil
		}
	}
	return nil, fmt.Errorf("invalid language: %s", enumString)
}
