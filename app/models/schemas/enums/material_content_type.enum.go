package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

/* ============================== Definition ============================== */

type MaterialContentType string

const (
	// basic types
	MaterialContentType_PlainText MaterialContentType = "text/plain"
	MaterialContentType_HTML      MaterialContentType = "text/html"
	MaterialContentType_PNG       MaterialContentType = "image/png"
	MaterialContentType_JPG       MaterialContentType = "image/jpg"
	MaterialContentType_JPEG      MaterialContentType = "image/jpeg"
	MaterialContentType_GIF       MaterialContentType = "image/gif"
	MaterialContentType_SVG       MaterialContentType = "image/svg"

	// custom types of Notezy
	// some charts, cards, drawing boards, etc.
)

/* ============================== All Instances ============================== */

var AllMaterialContentTypes = []MaterialContentType{
	MaterialContentType_PlainText,
	MaterialContentType_HTML,
	MaterialContentType_PNG,
	MaterialContentType_JPG,
	MaterialContentType_JPEG,
	MaterialContentType_GIF,
	MaterialContentType_SVG,
}

var AllMaterialContentTypeStrings = []string{
	string(MaterialContentType_PlainText),
	string(MaterialContentType_HTML),
	string(MaterialContentType_PNG),
	string(MaterialContentType_JPG),
	string(MaterialContentType_JPEG),
	string(MaterialContentType_GIF),
	string(MaterialContentType_SVG),
}

/* ============================== Methods ============================== */

func (mct MaterialContentType) Name() string {
	return reflect.TypeOf(mct).Name()
}

func (mct *MaterialContentType) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*mct = MaterialContentType(string(v))
		return nil
	case string:
		*mct = MaterialContentType(v)
		return nil
	}
	return scanError(value, mct)
}

func (mct MaterialContentType) Value() (driver.Value, error) {
	return string(mct), nil
}

func (mct MaterialContentType) String() string {
	return string(mct)
}

func (mct *MaterialContentType) IsValidEnum() bool {
	return slices.Contains(AllMaterialContentTypes, *mct)
}

func ConvertStringToMaterialContentType(enumString string) (*MaterialContentType, error) {
	for _, materialContentType := range AllMaterialContentTypes {
		if string(materialContentType) == enumString {
			return &materialContentType, nil
		}
	}
	return nil, fmt.Errorf("invalid material content type: %s", enumString)
}
