package blocknote

import "encoding/json"

/* ==================== InlineContent  Definitions ==================== */

type InlineContentType string

const InlineContentType_StyledText InlineContentType = "text"
const InlineContentType_Link InlineContentType = "link"

// place your custom inline content type above with the type of InlineContentType

type StyledText struct {
	Type   InlineContentType `json:"type" validate:"required,eq=text"`
	Text   string            `json:"text" validate:"required"`
	Styles []string          `json:"styles" validate:"omitempty"`
}

type Link struct {
	Type    InlineContentType `json:"type" validate:"required,eq=link"`
	Content []StyledText      `json:"content" validate:"omitempty"`
	Href    string            `json:"href" validate:"required"`
}

type CustomInlineContent struct {
	Type    InlineContentType      `json:"type" validate:"required"`
	Content []StyledText           `json:"content" validate:"omitempty"`
	Props   map[string]interface{} `json:"props" validate:"omitempty"`
}

// InlineContent = StyledText | Link | CustomInlineContent
type InlineContentUnion interface{ isInlineContent() bool }

func (*StyledText) isInlineContent() bool          { return true }
func (*Link) isInlineContent() bool                { return true }
func (*CustomInlineContent) isInlineContent() bool { return true }

type InlineContent struct {
	InlineContentUnion
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
		ic.InlineContentUnion = &v
	case InlineContentType_Link:
		var v Link
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		ic.InlineContentUnion = &v
	default:
		var v CustomInlineContent
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		ic.InlineContentUnion = &v
	}
	return nil
}

func (ic *InlineContent) MarshalJSON() ([]byte, error) {
	switch v := ic.InlineContentUnion.(type) {
	case *StyledText:
		return json.Marshal(v)
	case *Link:
		return json.Marshal(v)
	case *CustomInlineContent:
		return json.Marshal(v)
	default:
		return json.Marshal(nil)
	}
}

func NewStyledText(text string, styles []string) *StyledText {
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

func NewCustomInlineContent(customType InlineContentType, content []StyledText, props map[string]interface{}) *CustomInlineContent {
	if props == nil {
		props = map[string]interface{}{}
	}
	return &CustomInlineContent{
		Type:    customType,
		Content: content,
		Props:   props,
	}
}
