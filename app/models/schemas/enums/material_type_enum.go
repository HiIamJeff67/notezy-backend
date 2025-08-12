package enums

import (
	"database/sql/driver"
	"reflect"
	"slices"
)

/* ============================== Definition ============================== */

type MaterialType string

const (
	MaterialType_Textbook      MaterialType = "Textbook"
	MaterialType_Notebook      MaterialType = "Notebook"
	MaterialType_LearningCards MaterialType = "LearningCards"
)

/* ============================== All Instances ============================== */

var AllMaterialTypes = []MaterialType{
	MaterialType_Textbook,
	MaterialType_Notebook,
	MaterialType_LearningCards,
}

var AllMaterialTypeStrings = []string{
	string(MaterialType_Textbook),
	string(MaterialType_Notebook),
	string(MaterialType_LearningCards),
}

/* ============================== Methods ============================== */

func (m MaterialType) Name() string {
	return reflect.TypeOf(m).Name()
}

func (m *MaterialType) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*m = MaterialType(string(v))
		return nil
	case string:
		*m = MaterialType(v)
		return nil
	}
	return scanError(value, m)
}

func (m MaterialType) Value() (driver.Value, error) {
	return string(m), nil
}

func (m MaterialType) String() string {
	return string(m)
}

func (m *MaterialType) IsValidEnum() bool {
	return slices.Contains(AllMaterialTypes, *m)
}
