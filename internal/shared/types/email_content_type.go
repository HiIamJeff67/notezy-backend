package types

import "slices"

type EmailContentType string

const (
	EmailContentType_PlainText EmailContentType = "text/plain"
	EmailContentType_HTML      EmailContentType = "text/html"
	EmailContentType_Markdown  EmailContentType = "text/markdown"
)

func (ct EmailContentType) String() string {
	return string(ct)
}

func (ct *EmailContentType) IsValidEnum() bool {
	return slices.Contains(AllEmailContentTypes, *ct)
}

/* ========================= All EmailContentTypes ========================= */
var AllEmailContentTypes = []EmailContentType{
	EmailContentType_PlainText,
	EmailContentType_HTML,
	EmailContentType_Markdown,
}
var AllEmailContentTypeStrings = []string{
	string(EmailContentType_PlainText),
	string(EmailContentType_HTML),
	string(EmailContentType_Markdown),
}
