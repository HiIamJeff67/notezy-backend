package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

/* ============================== Definition ============================== */

type UsersToBillingPlansStatus string

const (
	UsersToBillingPlansStatus_ApprovalPending = "APPROVAL_PENDING"
	UsersToBillingPlansStatus_Approved        = "APPROVED"
	UsersToBillingPlansStatus_Active          = "ACTIVE"
	UsersToBillingPlansStatus_Suspended       = "SUSPENDED"
	UsersToBillingPlansStatus_Cancelled       = "CANCELLED"
	UsersToBillingPlansStatus_Expired         = "EXPIRED"
)

/* ============================== All Instances ============================== */

var AllUsersToBillingPlansStatuses = []UsersToBillingPlansStatus{
	UsersToBillingPlansStatus_ApprovalPending,
	UsersToBillingPlansStatus_Approved,
	UsersToBillingPlansStatus_Active,
	UsersToBillingPlansStatus_Suspended,
	UsersToBillingPlansStatus_Cancelled,
	UsersToBillingPlansStatus_Expired,
}
var AllUsersToBillingPlansStatusStrings = []string{
	string(UsersToBillingPlansStatus_ApprovalPending),
	string(UsersToBillingPlansStatus_Approved),
	string(UsersToBillingPlansStatus_Active),
	string(UsersToBillingPlansStatus_Suspended),
	string(UsersToBillingPlansStatus_Cancelled),
	string(UsersToBillingPlansStatus_Expired),
}

/* ============================== Methods ============================== */

func (utbps UsersToBillingPlansStatus) Name() string {
	return reflect.TypeOf(utbps).Name()
}

func (utbps *UsersToBillingPlansStatus) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*utbps = UsersToBillingPlansStatus(string(v))
		return nil
	case string:
		*utbps = UsersToBillingPlansStatus(v)
		return nil
	}
	return scanError(value, utbps)
}

func (utbps UsersToBillingPlansStatus) Value() (driver.Value, error) {
	return string(utbps), nil
}

func (utbps UsersToBillingPlansStatus) String() string {
	return string(utbps)
}

func (utbps *UsersToBillingPlansStatus) IsValidEnum() bool {
	for _, enum := range AllUsersToBillingPlansStatuses {
		if *utbps == enum {
			return true
		}
	}
	return false
}

func ConvertStringToUsersToBillingPlansStatus(enumString string) (*UsersToBillingPlansStatus, error) {
	for _, supportedCurrencyCode := range AllUsersToBillingPlansStatuses {
		if string(supportedCurrencyCode) == enumString {
			return &supportedCurrencyCode, nil
		}
	}
	return nil, fmt.Errorf("invalid users to billing plans status: %s", enumString)
}
