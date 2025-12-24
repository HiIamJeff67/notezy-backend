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
