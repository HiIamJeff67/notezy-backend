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

func (mt MaterialType) Name() string {
	return reflect.TypeOf(mt).Name()
}

func (mt *MaterialType) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*mt = MaterialType(string(v))
		return nil
	case string:
		*mt = MaterialType(v)
		return nil
	}
	return scanError(value, mt)
}

func (mt MaterialType) Value() (driver.Value, error) {
	return string(mt), nil
}

func (mt MaterialType) String() string {
	return string(mt)
}

func (mt *MaterialType) IsValidEnum() bool {
	return slices.Contains(AllMaterialTypes, *mt)
}
