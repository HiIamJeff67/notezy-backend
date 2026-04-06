package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

/* ============================== Definition ============================== */

type BillingIntervalUnit string

const (
	BillingIntervalUnit_Day   BillingIntervalUnit = "DAY"
	BillingIntervalUnit_Week  BillingIntervalUnit = "WEEK"
	BillingIntervalUnit_Month BillingIntervalUnit = "MONTH"
	BillingIntervalUnit_Year  BillingIntervalUnit = "YEAR"
)

/* ============================== All Instances ============================== */

var AllBillingIntervalUnits = []BillingIntervalUnit{
	BillingIntervalUnit_Day,
	BillingIntervalUnit_Week,
	BillingIntervalUnit_Month,
	BillingIntervalUnit_Year,
}
var AllBillingIntervalUnitStrings = []string{
	string(BillingIntervalUnit_Day),
	string(BillingIntervalUnit_Week),
	string(BillingIntervalUnit_Month),
	string(BillingIntervalUnit_Year),
}

/* ============================== Methods ============================== */

func (biu BillingIntervalUnit) Name() string {
	return reflect.TypeOf(biu).Name()
}

func (biu *BillingIntervalUnit) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*biu = BillingIntervalUnit(string(v))
		return nil
	case string:
		*biu = BillingIntervalUnit(v)
		return nil
	}
	return scanError(value, biu)
}

func (biu BillingIntervalUnit) Value() (driver.Value, error) {
	return string(biu), nil
}

func (biu BillingIntervalUnit) String() string {
	return string(biu)
}

func (biu *BillingIntervalUnit) IsValidEnum() bool {
	for _, enum := range AllBillingIntervalUnits {
		if *biu == enum {
			return true
		}
	}
	return false
}

func ConvertStringToBillingIntervalUnit(enumString string) (*BillingIntervalUnit, error) {
	for _, supportedCurrencyCode := range AllBillingIntervalUnits {
		if string(supportedCurrencyCode) == enumString {
			return &supportedCurrencyCode, nil
		}
	}
	return nil, fmt.Errorf("invalid billing interval unit: %s", enumString)
}
