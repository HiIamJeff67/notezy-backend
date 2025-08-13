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
	MaterialContentType_PlainText MaterialContentType = "text/plain"
	MaterialContentType_HTML      MaterialContentType = "text/html"
	MaterialContentType_Markdown  MaterialContentType = "text/markdown"
	MaterialContentType_PNG       MaterialContentType = "image/png"
	MaterialContentType_JPG       MaterialContentType = "image/jpg"
	MaterialContentType_JPEG      MaterialContentType = "image/jpeg"
	MaterialContentType_GIF       MaterialContentType = "image/gif"
	MaterialContentType_SVG       MaterialContentType = "image/svg"
	MaterialContentType_MP3       MaterialContentType = "video/mp3"
	MaterialContentType_MP4       MaterialContentType = "video/mp4"

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
