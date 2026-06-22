package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

type RoutinePeriod string

const (
	RoutinePeriod_Daily   RoutinePeriod = "Daily"
	RoutinePeriod_Weekly  RoutinePeriod = "Weekly"
	RoutinePeriod_Monthly RoutinePeriod = "Monthly"
)

var AllRoutinePeriods = []RoutinePeriod{
	RoutinePeriod_Daily,
	RoutinePeriod_Weekly,
	RoutinePeriod_Monthly,
}

var AllRoutinePeriodStrings = []string{
	string(RoutinePeriod_Daily),
	string(RoutinePeriod_Weekly),
	string(RoutinePeriod_Monthly),
}

func (rp RoutinePeriod) Name() string {
	return reflect.TypeOf(rp).Name()
}

func (rp *RoutinePeriod) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*rp = RoutinePeriod(string(v))
		return nil
	case string:
		*rp = RoutinePeriod(v)
		return nil
	}
	return scanError(value, rp)
}

func (rp RoutinePeriod) Value() (driver.Value, error) {
	return string(rp), nil
}

func (rp RoutinePeriod) String() string {
	return string(rp)
}

func (rp *RoutinePeriod) IsValidEnum() bool {
	return slices.Contains(AllRoutinePeriods, *rp)
}

func ConvertStringToRoutinePeriod(enumString string) (*RoutinePeriod, error) {
	for _, routinePeriod := range AllRoutinePeriods {
		if string(routinePeriod) == enumString {
			return &routinePeriod, nil
		}
	}
	return nil, fmt.Errorf("invalid routine status: %s", enumString)
}
