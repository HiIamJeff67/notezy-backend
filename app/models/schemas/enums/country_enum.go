package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

/* ============================== Definition ============================== */

type Country string

const (
	Country_Taiwan                Country = "Taiwan"
	Country_Japan                 Country = "Japan"
	Country_Malaysia              Country = "Malaysia"
	Country_Singapore             Country = "Singapore"
	Country_China                 Country = "China"
	Country_UnitedStatusOfAmerica Country = "UnitedStatesOfAmerica"
	Country_UnitedKingdom         Country = "UnitedKingdom"
	Country_Australia             Country = "Australia"
	Country_Canada                Country = "Canada"
)

/* ========================= All Instances ========================= */

var AllCountries = []Country{
	Country_Taiwan,
	Country_Japan,
	Country_Malaysia,
	Country_Singapore,
	Country_China,
	Country_UnitedStatusOfAmerica,
	Country_UnitedKingdom,
	Country_Australia,
	Country_Canada,
}
var AllCountryStrings = []string{
	string(Country_Taiwan),
	string(Country_Japan),
	string(Country_Malaysia),
	string(Country_Singapore),
	string(Country_China),
	string(Country_UnitedStatusOfAmerica),
	string(Country_UnitedKingdom),
	string(Country_Australia),
	string(Country_Canada),
}

/* ============================== Methods ============================== */

func (c Country) Name() string {
	return reflect.TypeOf(c).Name()
}

func (c *Country) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*c = Country(string(v))
		return nil
	case string:
		*c = Country(v)
		return nil
	}
	return scanError(value, c)
}

func (c Country) Value() (driver.Value, error) {
	return string(c), nil
}

func (c Country) String() string {
	return string(c)
}

func (c *Country) IsValidEnum() bool {
	for _, enum := range AllCountries {
		if *c == enum {
			return true
		}
	}
	return false
}

func ConvertStringToCountry(enumString string) (*Country, error) {
	for _, country := range AllCountries {
		if string(country) == enumString {
			return &country, nil
		}
	}
	return nil, fmt.Errorf("invalid country: %s", enumString)
}
