package types

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

func (u *UUIDArray) Scan(value interface{}) error {
	if value == nil {
		*u = UUIDArray{}
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into UUIDArray", value)
	}
	str = strings.Trim(str, "{}")
	if str == "" {
		*u = UUIDArray{}
		return nil
	}
	parts := strings.Split(str, ",")
	uuids := make(UUIDArray, len(parts))
	for i, part := range parts {
		id, err := uuid.Parse(strings.TrimSpace(part))
		if err != nil {
			return fmt.Errorf("invalid UUID in array: %s", part)
		}
		uuids[i] = id
	}
	*u = uuids
	return nil
}

func (u UUIDArray) Value() (driver.Value, error) {
	if len(u) == 0 {
		return "{}", nil
	}
	strs := make([]string, len(u))
	for i, id := range u {
		strs[i] = id.String()
	}
	return "{" + strings.Join(strs, ",") + "}", nil
}
