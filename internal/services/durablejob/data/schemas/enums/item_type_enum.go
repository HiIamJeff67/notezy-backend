package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

type ItemType string

const (
	ItemType_BlockPack ItemType = "BlockPack"
	ItemType_Material  ItemType = "Material"
)

var AllItemTypes = []ItemType{
	ItemType_BlockPack,
	ItemType_Material,
}

var AllItemTypeStrings = []string{
	string(ItemType_BlockPack),
	string(ItemType_Material),
}

func (it ItemType) Name() string {
	return reflect.TypeOf(it).Name()
}

func (it *ItemType) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*it = ItemType(string(v))
		return nil
	case string:
		*it = ItemType(v)
		return nil
	}
	return scanError(value, it)
}

func (it ItemType) Value() (driver.Value, error) {
	return string(it), nil
}

func (it ItemType) String() string {
	return string(it)
}

func (it *ItemType) IsValidEnum() bool {
	return slices.Contains(AllItemTypes, *it)
}

func ConvertStringToItemType(enumString string) (*ItemType, error) {
	for _, itemType := range AllItemTypes {
		if string(itemType) == enumString {
			return &itemType, nil
		}
	}
	return nil, fmt.Errorf("invalid item type: %s", enumString)
}
