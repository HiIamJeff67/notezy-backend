package enums

import (
	"database/sql/driver"
	"reflect"
	"slices"
)

/* ============================== Definition ============================== */

type MaterialContentType string

const (
	// basic types
	MaterialContentType_PlainText MaterialContentType = "Text_Plain"
	MaterialContentType_HTML      MaterialContentType = "Text_HTML"
	MaterialContentType_Markdown  MaterialContentType = "Text_Markdown"
	MaterialContentType_PNG       MaterialContentType = "Image_PNG"
	MaterialContentType_JPG       MaterialContentType = "Image_JPG"
	MaterialContentType_JPEG      MaterialContentType = "Image_JPEG"
	MaterialContentType_GIF       MaterialContentType = "Image_GIF"
	MaterialContentType_SVG       MaterialContentType = "Image_SVG"
	MaterialContentType_MP3       MaterialContentType = "Video_MP3"
	MaterialContentType_MP4       MaterialContentType = "Video_MP4"

	// custom types of Notezy
	// some charts, cards, drawing boards, etc.
)

/* ============================== All Instances ============================== */

var AllMaterialContentTypes = []MaterialContentType{
	MaterialContentType_PlainText,
	MaterialContentType_HTML,
	MaterialContentType_Markdown,
	MaterialContentType_PNG,
	MaterialContentType_JPG,
	MaterialContentType_JPEG,
	MaterialContentType_GIF,
	MaterialContentType_SVG,
	MaterialContentType_MP3,
	MaterialContentType_MP4,
}

var AllMaterialContentTypeStrings = []string{
	string(MaterialContentType_PlainText),
	string(MaterialContentType_HTML),
	string(MaterialContentType_Markdown),
	string(MaterialContentType_PNG),
	string(MaterialContentType_JPG),
	string(MaterialContentType_JPEG),
	string(MaterialContentType_GIF),
	string(MaterialContentType_SVG),
	string(MaterialContentType_MP3),
	string(MaterialContentType_MP4),
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
