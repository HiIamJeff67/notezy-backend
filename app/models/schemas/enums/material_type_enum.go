package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

/* ============================== Definition ============================== */

type MaterialType string

const (
	MaterialType_Textbook      MaterialType = "Textbook"      // BlockNote(HTML or Markdown)
	MaterialType_Notebook      MaterialType = "Notebook"      // BlockNote(Markdown)
	MaterialType_LearningCards MaterialType = "LearningCards" // BlockNote(HTML)
	MaterialType_Workflow      MaterialType = "Workflow"      // ReactFlow(Canva)
)

// the sub type to indicate the content type of material files
type MaterialContentType string

const (
	// basic types
	MaterialContentType_JSON      MaterialContentType = "application/json"
	MaterialContentType_PlainText MaterialContentType = "text/plain"
	MaterialContentType_HTML      MaterialContentType = "text/html"
	MaterialContentType_PNG       MaterialContentType = "image/png"
	MaterialContentType_JPG       MaterialContentType = "image/jpg"
	MaterialContentType_JPEG      MaterialContentType = "image/jpeg"
	MaterialContentType_GIF       MaterialContentType = "image/gif"
	MaterialContentType_SVG       MaterialContentType = "image/svg"
)

/* ============================== All Instances ============================== */

var AllMaterialTypes = []MaterialType{
	MaterialType_Textbook,
	MaterialType_Notebook,
	MaterialType_LearningCards,
	MaterialType_Workflow,
}

var AllMaterialTypeStrings = []string{
	string(MaterialType_Textbook),
	string(MaterialType_Notebook),
	string(MaterialType_LearningCards),
	string(MaterialType_Workflow),
}

// the sub types array for the content type of material files
var AllMaterialContentTypes = []MaterialContentType{
	MaterialContentType_PlainText,
	MaterialContentType_HTML,
	MaterialContentType_PNG,
	MaterialContentType_JPG,
	MaterialContentType_JPEG,
	MaterialContentType_GIF,
	MaterialContentType_SVG,
}

// the sub type strings array for the content type of material files
var AllMaterialContentTypeStrings = []string{
	string(MaterialContentType_PlainText),
	string(MaterialContentType_HTML),
	string(MaterialContentType_PNG),
	string(MaterialContentType_JPG),
	string(MaterialContentType_JPEG),
	string(MaterialContentType_GIF),
	string(MaterialContentType_SVG),
}

// mapping from material types to their allowed content types,
// this should be provided to the frontend and useless in the backend
var MaterialTypeToAllowedContentTypes = map[MaterialType][]MaterialContentType{
	MaterialType_Textbook: {
		MaterialContentType_HTML,
		MaterialContentType_PlainText,
	},
	MaterialType_Notebook: {
		MaterialContentType_HTML,
		MaterialContentType_PlainText,
	},
	MaterialType_LearningCards: {
		MaterialContentType_HTML,
	},
	MaterialType_Workflow: {},
}

// mapping from material types to their allowed content type strings,
// this should be provided to the frontend and useless in the backend
var MaterialTypeToAllowedContentTypeStrings = map[MaterialType][]string{
	MaterialType_Textbook: {
		MaterialContentType_HTML.String(),
		MaterialContentType_PlainText.String(),
	},
	MaterialType_Notebook: {
		MaterialContentType_HTML.String(),
		MaterialContentType_PlainText.String(),
	},
	MaterialType_LearningCards: {
		MaterialContentType_HTML.String(),
	},
	MaterialType_Workflow: {},
}

/* ============================== Methods for MaterialType ============================== */

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

func (mt MaterialType) IsContentTypeAllowed(contentType MaterialContentType) bool {
	for _, allowedContentType := range MaterialTypeToAllowedContentTypes[mt] {
		if contentType == allowedContentType {
			return true
		}
	}

	return false
}

func (mt MaterialType) IsContentTypeStringAllowed(contentTypeString string) bool {
	for _, allowedContentTypeString := range MaterialTypeToAllowedContentTypeStrings[mt] {
		if strings.Contains(allowedContentTypeString, contentTypeString) {
			return true
		}
	}

	return false
}

func (mt MaterialType) AllowedContentTypes() []MaterialContentType {
	return MaterialTypeToAllowedContentTypes[mt]
}

func (mt MaterialType) AllowedContentTypeStrings() []string {
	return MaterialTypeToAllowedContentTypeStrings[mt]
}

func ConvertStringToMaterialType(enumString string) (*MaterialType, error) {
	for _, materialType := range AllMaterialTypes {
		if string(materialType) == enumString {
			return &materialType, nil
		}
	}
	return nil, fmt.Errorf("invalid material type: %s", enumString)
}

/* ============================== Methods for MaterialContentType ============================== */

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
