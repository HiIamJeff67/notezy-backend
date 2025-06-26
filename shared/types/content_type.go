package types

import "slices"

type ContentType string

const (
	ContentType_PlainText ContentType = "text/plain"
	ContentType_HTML      ContentType = "text/html"
	ContentType_Markdown  ContentType = "text/markdown"
)

func (ct ContentType) String() string {
	return string(ct)
}

func (ct *ContentType) IsValidEnum() bool {
	return slices.Contains(AllContentTypes, *ct)
}

/* ========================= All ContentTypes ========================= */
var AllContentTypes = []ContentType{
	ContentType_PlainText,
	ContentType_HTML,
	ContentType_Markdown,
}
var AllContentTypeStrings = []string{
	string(ContentType_PlainText),
	string(ContentType_HTML),
	string(ContentType_Markdown),
}
