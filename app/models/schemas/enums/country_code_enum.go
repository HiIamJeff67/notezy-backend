package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

/* ============================== Definition ============================== */

type CountryCode string

const (
	CountryCode_Taiwan        CountryCode = "COUNTRY_CODE_886"
	CountryCode_Japan         CountryCode = "COUNTRY_CODE_81"
	CountryCode_Malaysia      CountryCode = "COUNTRY_CODE_60"
	CountryCode_Singapore     CountryCode = "COUNTRY_CODE_65"
	CountryCode_China         CountryCode = "COUNTRY_CODE_86"
	CountryCode_NANP          CountryCode = "COUNTRY_CODE_1"
	CountryCode_UnitedKingdom CountryCode = "COUNTRY_CODE_44"
	CountryCode_Australia     CountryCode = "COUNTRY_CODE_61"
)

/* ========================= All Instances ========================= */

var AllCountryCodes = []CountryCode{
	CountryCode_Taiwan,
	CountryCode_Japan,
	CountryCode_Malaysia,
	CountryCode_Singapore,
	CountryCode_China,
	CountryCode_NANP, // NANP stands for North American Numbering Plan, it's used in United States of America and Canada
	CountryCode_UnitedKingdom,
	CountryCode_Australia,
}
var AllCountryCodeStrings = []string{
	string(CountryCode_Taiwan),
	string(CountryCode_Japan),
	string(CountryCode_Malaysia),
	string(CountryCode_Singapore),
	string(CountryCode_China),
	string(CountryCode_NANP),
	string(CountryCode_UnitedKingdom),
	string(CountryCode_Australia),
}

/* ============================== Methods ============================== */

func (cc CountryCode) Name() string {
	return reflect.TypeOf(cc).Name()
}

func (cc *CountryCode) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*cc = CountryCode(string(v))
		return nil
	case string:
		*cc = CountryCode(v)
		return nil
	}
	return scanError(value, cc)
}

func (cc CountryCode) Value() (driver.Value, error) {
	return string(cc), nil
}

func (cc CountryCode) String() string {
	return string(cc)
}

func (cc *CountryCode) IsValidEnum() bool {
	for _, enum := range AllCountryCodes {
		if *cc == enum {
			return true
		}
	}
	return false
}

func ConvertStringToCountryCode(enumString string) (*CountryCode, error) {
	for _, countryCode := range AllCountryCodes {
		if string(countryCode) == enumString {
			return &countryCode, nil
		}
	}
	return nil, fmt.Errorf("invalid country code: %s", enumString)
}
