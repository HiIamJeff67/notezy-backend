package enums

import (
	"database/sql/driver"
	"reflect"
)

/* ============================== CountryCode Definition ============================== */
type CountryCode string

const (
	CountryCode_Default       CountryCode = "Default" // null value
	CountryCode_Taiwan        CountryCode = "+886"
	CountryCode_Japan         CountryCode = "+81"
	CountryCode_Malaysia      CountryCode = "+60"
	CountryCode_Singapore     CountryCode = "+65"
	CountryCode_China         CountryCode = "+86"
	CountryCode_NANP          CountryCode = "+1"
	CountryCode_UnitedKingdom CountryCode = "+44"
	CountryCode_Australia     CountryCode = "+61"
)

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

func (cc *CountryCode) IsValidEnum() bool {
	for _, enum := range AllCountryCodes {
		if *cc == enum {
			return true
		}
	}
	return false
}

/* ========================= All CountryCodes ========================= */
var AllCountryCodes = []CountryCode{
	CountryCode_Default,
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
	string(CountryCode_Default),
	string(CountryCode_Taiwan),
	string(CountryCode_Japan),
	string(CountryCode_Malaysia),
	string(CountryCode_Singapore),
	string(CountryCode_China),
	string(CountryCode_NANP),
	string(CountryCode_UnitedKingdom),
	string(CountryCode_Australia),
}
