package types

type ContentType string

const (
	ContentType_PlainText ContentType = "text/plain"
	ContentType_HTML      ContentType = "text/html"
	ContentType_Markdown  ContentType = "text/markdown"
)

func (ct ContentType) Value() (string, error) {
	return string(ct), nil
}

func (ct *ContentType) IsValidEnum() bool {
	for _, enum := range AllContentTypes {
		if *ct == enum {
			return true
		}
	}
	return false
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
