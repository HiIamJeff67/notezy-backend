package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
)

type BillingPlanStatus enumcontract.BillingPlanStatus

func (value *BillingPlanStatus) ToContractable() *enumcontract.BillingPlanStatus {
	if value == nil {
		return nil
	}

	contractValue := enumcontract.BillingPlanStatus(*value)
	return &contractValue
}

func (value *BillingPlanStatus) ToStorable() *BillingPlanStatus {
	if value == nil {
		return nil
	}

	storableValue := *value
	return &storableValue
}

const (
	BillingPlanStatus_Created  BillingPlanStatus = BillingPlanStatus(enumcontract.BillingPlanStatus_Created)
	BillingPlanStatus_Active   BillingPlanStatus = BillingPlanStatus(enumcontract.BillingPlanStatus_Active)
	BillingPlanStatus_Inactive BillingPlanStatus = BillingPlanStatus(enumcontract.BillingPlanStatus_Inactive)
)

var AllBillingPlanStatuses = []BillingPlanStatus{
	BillingPlanStatus_Created,
	BillingPlanStatus_Active,
	BillingPlanStatus_Inactive,
}
var AllBillingPlanStatusStrings = []string{
	string(BillingPlanStatus_Created),
	string(BillingPlanStatus_Active),
	string(BillingPlanStatus_Inactive),
}

func (bps BillingPlanStatus) Name() string {
	return reflect.TypeOf(bps).Name()
}

func (bps *BillingPlanStatus) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*bps = BillingPlanStatus(string(v))
		return nil
	case string:
		*bps = BillingPlanStatus(v)
		return nil
	}
	return scanError(value, bps)
}

func (bps BillingPlanStatus) Value() (driver.Value, error) {
	return string(bps), nil
}

func (bps BillingPlanStatus) String() string {
	return string(bps)
}

func (bps *BillingPlanStatus) IsValidEnum() bool {
	for _, enum := range AllBillingPlanStatuses {
		if *bps == enum {
			return true
		}
	}
	return false
}

func ConvertStringToBillingPlanStatus(enumString string) (*BillingPlanStatus, error) {
	for _, supportedCurrencyCode := range AllBillingPlanStatuses {
		if string(supportedCurrencyCode) == enumString {
			return &supportedCurrencyCode, nil
		}
	}
	return nil, fmt.Errorf("invalid billing plan status: %s", enumString)
}
