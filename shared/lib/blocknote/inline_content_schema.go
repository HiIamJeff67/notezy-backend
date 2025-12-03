package blocknote

import (
	"encoding/json"

	validation "notezy-backend/app/validation"
)

/* ============================== Other Type Definition ============================== */

// place your custom inline content type above with the type of InlineContentType and implement some methods for it below

type Styles struct {
	Bold            bool   `json:"bold,omitempty" validate:"omitempty"`
	Italic          bool   `json:"italic,omitempty" validate:"omitempty"`
	Underline       bool   `json:"underline,omitempty" validate:"omitempty"`
	Strike          bool   `json:"strike,omitempty" validate:"omitempty"`
	Code            bool   `json:"code,omitempty" validate:"omitempty"`
	TextColor       string `json:"textColor,omitempty" validate:"omitempty,ishexcode"`
	BackgroundColor string `json:"backgroundColor,omitempty" validate:"omitempty,ishexcode"`
}

/* ============================== InlineContentUnion ============================== */

// InlineContent = Styles | StyledText | Link | CustomInlineContent
type InlineContentUnion interface {
	isInlineContent() bool
	Validate() error
}

type InlineContentType string

const InlineContentType_StyledText InlineContentType = "text"
const InlineContentType_Link InlineContentType = "link"

type StyledText struct {
	Type   InlineContentType `json:"type" validate:"required,eq=text"`
	Text   string            `json:"text" validate:"required,max=4096"`
	Styles Styles            `json:"styles" validate:"omitempty"`
}

func (*StyledText) isInlineContent() bool { return true }

func (st *StyledText) Validate() error { return validation.Validator.Struct(st) }

type Link struct {
	Type    InlineContentType `json:"type" validate:"required,eq=link"`
	Href    string            `json:"href" validate:"required,isurl"`
	Content []StyledText      `json:"content" validate:"omitempty,dive"` // use dive to validate recursively
}

func (*Link) isInlineContent() bool { return true }

func (l *Link) Validate() error { return validation.Validator.Struct(l) }

// type CustomInlineContent struct {
// 	Type    InlineContentType      `json:"type" validate:"required"`
// 	Props   map[string]interface{} `json:"props" validate:"omitempty"`
// 	Content []StyledText           `json:"content" validate:"omitempty,dive"`
// }

// func (*CustomInlineContent) isInlineContent() bool { return true }

// func (cic *CustomInlineContent) Validate() error { // calling the validator to validate the struct of cic }

/* ============================== InlineContent ============================== */

type InlineContent struct {
	InlineContentUnion
}

func (ic *InlineContent) Validate() error {
	if ic.InlineContentUnion == nil {
		return nil
	}
	return ic.InlineContentUnion.Validate()
}

func (ic *InlineContent) UnmarshalJSON(b []byte) error {
	var t struct {
		Type InlineContentType `json:"type" validate:"required"`
	}
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	switch t.Type {
	case InlineContentType_StyledText:
		var v StyledText
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		if err := v.Validate(); err != nil {
			return err
		}
		ic.InlineContentUnion = &v
	case InlineContentType_Link:
		var v Link
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		if err := v.Validate(); err != nil {
			return err
		}
		ic.InlineContentUnion = &v
		// case InlineContentType_CustomInlineContent:
		// 	var v CustomInlineContent
		// 	if err := json.Unmarshal(b, &v); err != nil {
		// 		return err
		// 	}
		// 	ic.InlineContentUnion = &v
	}
	return nil
}

func (ic *InlineContent) MarshalJSON() ([]byte, error) {
	switch v := ic.InlineContentUnion.(type) {
	case *StyledText:
		return json.Marshal(v)
	case *Link:
		return json.Marshal(v)
	// case *CustomInlineContent:
	// 	return json.Marshal(v)
	default:
		return json.Marshal(nil)
	}
}

func NewStyledText(text string, styles Styles) *StyledText {
	return &StyledText{
		Type:   InlineContentType_StyledText,
		Text:   text,
		Styles: styles,
	}
}

func NewLink(href string, content []StyledText) *Link {
	return &Link{
		Type:    InlineContentType_Link,
		Href:    href,
		Content: content,
	}
}

// func NewCustomInlineContent(customType InlineContentType, content []StyledText, props map[string]interface{}) *CustomInlineContent {
// 	if props == nil {
// 		props = map[string]interface{}{}
// 	}
// 	return &CustomInlineContent{
// 		Type:    customType,
// 		Content: content,
// 		Props:   props,
// 	}
// }
