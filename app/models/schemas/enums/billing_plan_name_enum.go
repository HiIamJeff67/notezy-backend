package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type BillingPlanName string

const (
	BillingPlanName_NotezyMonthlyFreePlan       BillingPlanName = "Notezy Monthly Free Plan"
	BillingPlanName_NotezyMonthlyProPlan        BillingPlanName = "Notezy Monthly Pro Plan"
	BillingPlanName_NotezyYearlyProPlan         BillingPlanName = "Notezy Yearly Pro Plan"
	BillingPlanName_NotezyMonthlyPremiumPlan    BillingPlanName = "Notezy Monthly Premium Plan"
	BillingPlanName_NotezyYearlyPremiumPlan     BillingPlanName = "Notezy Yearly Premium Plan"
	BillingPlanName_NotezyMonthlyUltimatePlan   BillingPlanName = "Notezy Monthly Ultimate Plan"
	BillingPlanName_NotezyYearlyUltimatePlan    BillingPlanName = "Notezy Yearly Ultimate Plan"
	BillingPlanName_NotezyMonthlyEnterprisePlan BillingPlanName = "Notezy Monthly Enterprise Plan"
	BillingPlanName_NotezyYearlyEnterprisePlan  BillingPlanName = "Notezy Yearly Enterprise Plan"
)

var AllBillingPlanNames = []BillingPlanName{
	BillingPlanName_NotezyMonthlyFreePlan,
	BillingPlanName_NotezyMonthlyProPlan,
	BillingPlanName_NotezyYearlyProPlan,
	BillingPlanName_NotezyMonthlyPremiumPlan,
	BillingPlanName_NotezyYearlyPremiumPlan,
	BillingPlanName_NotezyMonthlyUltimatePlan,
	BillingPlanName_NotezyYearlyUltimatePlan,
	BillingPlanName_NotezyMonthlyEnterprisePlan,
	BillingPlanName_NotezyYearlyEnterprisePlan,
}

var AllBillingPlanNameStrings = []string{
	string(BillingPlanName_NotezyMonthlyFreePlan),
	string(BillingPlanName_NotezyMonthlyProPlan),
	string(BillingPlanName_NotezyYearlyProPlan),
	string(BillingPlanName_NotezyMonthlyPremiumPlan),
	string(BillingPlanName_NotezyYearlyPremiumPlan),
	string(BillingPlanName_NotezyMonthlyUltimatePlan),
	string(BillingPlanName_NotezyYearlyUltimatePlan),
	string(BillingPlanName_NotezyMonthlyEnterprisePlan),
	string(BillingPlanName_NotezyYearlyEnterprisePlan),
}

func (bpn BillingPlanName) Name() string {
	return reflect.TypeOf(bpn).Name()
}

func (bpn *BillingPlanName) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*bpn = BillingPlanName(string(v))
		return nil
	case string:
		*bpn = BillingPlanName(v)
		return nil
	}
	return scanError(value, bpn)
}

func (bpn BillingPlanName) Value() (driver.Value, error) {
	return string(bpn), nil
}

func (bpn BillingPlanName) String() string {
	return string(bpn)
}

func (bpn *BillingPlanName) IsValidEnum() bool {
	for _, enum := range AllBillingPlanNames {
		if *bpn == enum {
			return true
		}
	}
	return false
}

func ConvertStringToBillingPlanName(enumString string) (*BillingPlanName, error) {
	for _, billingPlanName := range AllBillingPlanNames {
		if string(billingPlanName) == enumString {
			return &billingPlanName, nil
		}
	}
	return nil, fmt.Errorf("invalid billing plan name: %s", enumString)
}
