package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type SupportedCurrencyCode string

const (
	SupportedCurrencyCode_USD SupportedCurrencyCode = "USD"
	SupportedCurrencyCode_EUR SupportedCurrencyCode = "EUR"
	SupportedCurrencyCode_JPY SupportedCurrencyCode = "JPY"
	SupportedCurrencyCode_TWD SupportedCurrencyCode = "TWD"
	SupportedCurrencyCode_KRW SupportedCurrencyCode = "KRW"
	SupportedCurrencyCode_CNY SupportedCurrencyCode = "CNY"
)

/* ============================== All instances ============================== */

var AllSupportedCurrencyCodes = []SupportedCurrencyCode{
	SupportedCurrencyCode_USD,
	SupportedCurrencyCode_EUR,
	SupportedCurrencyCode_JPY,
	SupportedCurrencyCode_TWD,
	SupportedCurrencyCode_KRW,
	SupportedCurrencyCode_CNY,
}
var AllSupportedCurrencyCodeStrings = []string{
	string(SupportedCurrencyCode_USD),
	string(SupportedCurrencyCode_EUR),
	string(SupportedCurrencyCode_JPY),
	string(SupportedCurrencyCode_TWD),
	string(SupportedCurrencyCode_KRW),
	string(SupportedCurrencyCode_CNY),
}

func (scc SupportedCurrencyCode) Name() string {
	return reflect.TypeOf(scc).Name()
}

func (scc *SupportedCurrencyCode) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*scc = SupportedCurrencyCode(string(v))
		return nil
	case string:
		*scc = SupportedCurrencyCode(v)
		return nil
	}
	return scanError(value, scc)
}

func (scc SupportedCurrencyCode) Value() (driver.Value, error) {
	return string(scc), nil
}

func (scc SupportedCurrencyCode) String() string {
	return string(scc)
}

func (scc *SupportedCurrencyCode) IsValidEnum() bool {
	for _, enum := range AllSupportedCurrencyCodes {
		if *scc == enum {
			return true
		}
	}
	return false
}

func ConvertStringToSupportedCurrencyCode(enumString string) (*SupportedCurrencyCode, error) {
	for _, supportedCurrencyCode := range AllSupportedCurrencyCodes {
		if string(supportedCurrencyCode) == enumString {
			return &supportedCurrencyCode, nil
		}
	}
	return nil, fmt.Errorf("invalid supported currency code: %s", enumString)
}
