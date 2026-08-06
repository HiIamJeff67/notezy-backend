package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
)

type BillingIntervalUnit enumcontract.BillingIntervalUnit

func (value *BillingIntervalUnit) ToContractable() *enumcontract.BillingIntervalUnit {
	if value == nil {
		return nil
	}

	contractValue := enumcontract.BillingIntervalUnit(*value)
	return &contractValue
}

func (value *BillingIntervalUnit) ToStorable() *BillingIntervalUnit {
	if value == nil {
		return nil
	}

	storableValue := *value
	return &storableValue
}

const (
	BillingIntervalUnit_Day   BillingIntervalUnit = BillingIntervalUnit(enumcontract.BillingIntervalUnit_Day)
	BillingIntervalUnit_Week  BillingIntervalUnit = BillingIntervalUnit(enumcontract.BillingIntervalUnit_Week)
	BillingIntervalUnit_Month BillingIntervalUnit = BillingIntervalUnit(enumcontract.BillingIntervalUnit_Month)
	BillingIntervalUnit_Year  BillingIntervalUnit = BillingIntervalUnit(enumcontract.BillingIntervalUnit_Year)
)

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
