package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

// MaterialContentType indicates the MIME content type of material files.
type MaterialContentType string

const (
	MaterialContentType_None      MaterialContentType = "none"
	MaterialContentType_JSON      MaterialContentType = "application/json"
	MaterialContentType_PDF       MaterialContentType = "application/pdf"
	MaterialContentType_PlainText MaterialContentType = "text/plain"
	MaterialContentType_HTML      MaterialContentType = "text/html"
	MaterialContentType_Markdown  MaterialContentType = "text/markdown"
	MaterialContentType_PNG       MaterialContentType = "image/png"
	MaterialContentType_JPG       MaterialContentType = "image/jpg"
	MaterialContentType_JPEG      MaterialContentType = "image/jpeg"
	MaterialContentType_GIF       MaterialContentType = "image/gif"
	MaterialContentType_SVG       MaterialContentType = "image/svg+xml"
	MaterialContentType_WebP      MaterialContentType = "image/webp"
	MaterialContentType_MP4       MaterialContentType = "video/mp4"
	MaterialContentType_WebM      MaterialContentType = "video/webm"
	MaterialContentType_Mpeg      MaterialContentType = "audio/mpeg"
)

var AllMaterialContentTypes = []MaterialContentType{
	MaterialContentType_None,
	MaterialContentType_JSON,
	MaterialContentType_PDF,
	MaterialContentType_PlainText,
	MaterialContentType_HTML,
	MaterialContentType_Markdown,
	MaterialContentType_PNG,
	MaterialContentType_JPG,
	MaterialContentType_JPEG,
	MaterialContentType_GIF,
	MaterialContentType_SVG,
	MaterialContentType_WebP,
	MaterialContentType_MP4,
	MaterialContentType_WebM,
	MaterialContentType_Mpeg,
}

var AllMaterialContentTypeStrings = []string{
	string(MaterialContentType_None),
	string(MaterialContentType_JSON),
	string(MaterialContentType_PDF),
	string(MaterialContentType_PlainText),
	string(MaterialContentType_HTML),
	string(MaterialContentType_Markdown),
	string(MaterialContentType_PNG),
	string(MaterialContentType_JPG),
	string(MaterialContentType_JPEG),
	string(MaterialContentType_GIF),
	string(MaterialContentType_SVG),
	string(MaterialContentType_WebP),
	string(MaterialContentType_MP4),
	string(MaterialContentType_WebM),
	string(MaterialContentType_Mpeg),
}

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
