package enums

import (
	"database/sql/driver"
	"reflect"
)

/* ============================== UserStatus Definition ============================== */
type UserStatus string

const (
	UserStatus_Online       UserStatus = "Online"
	UserStatus_AFK          UserStatus = "AFK"
	UserStatus_DoNotDisturb UserStatus = "DoNotDisturb"
	UserStatus_Offline      UserStatus = "Offline"
)

func (s *UserStatus) Name() string {
	return reflect.TypeOf(s).Name()
}

func (s *UserStatus) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*s = UserStatus(string(v))
		return nil
	case string:
		*s = UserStatus(v)
		return nil
	}
	return scanError(value, s)
}

func (s UserStatus) Value() (driver.Value, error) {
	return string(s), nil
}

var AllUserStatuses = []UserStatus{
	UserStatus_Online,
	UserStatus_AFK,
	UserStatus_DoNotDisturb,
	UserStatus_Offline,
}
var AllUserStatusStrings = []string{
	string(UserStatus_Online),
	string(UserStatus_AFK),
	string(UserStatus_DoNotDisturb),
	string(UserStatus_Offline),
}
