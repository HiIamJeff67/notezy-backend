package blocknote

import (
	"encoding/json"

	validation "notezy-backend/app/validation"
)

type BlockProps interface {
	IsBlockProps() bool
	Validate() error
}

/* ============================== BaseProps ============================== */

type BaseProps struct {
	TextColor       string `json:"textColor,omitempty" validate:"omitempty,iscolororhexcode"`
	BackgroundColor string `json:"backgroundColor,omitempty" validate:"omitempty,iscolororhexcode"`
	TextAlignment   string `json:"textAlignment,omitempty" validate:"omitempty,istextalignment"`
}

func (bp *BaseProps) IsBlockProps() bool { return true }

func (bp *BaseProps) Validate() error { return validation.Validator.Struct(bp) }

/* ============================== HeadingProps ============================== */

type HeadingProps struct {
	BaseProps
	Level        int  `json:"level" validate:"required,isheadinglevel"`
	IsToggleable bool `json:"isToggleable,omitempty" validate:"omitempty"`
}

func (hp *HeadingProps) IsBlockProps() bool { return true }

func (hp *HeadingProps) Validate() error { return validation.Validator.Struct(hp) }

/* ============================== CheckListItemProps ============================== */

type CheckListItemProps struct {
	BaseProps
	Checked bool `json:"checked,omitempty"`
}

func (clip *CheckListItemProps) IsBlockProps() bool { return true }

func (clip *CheckListItemProps) Validate() error { return validation.Validator.Struct(clip) }

/* ============================== FileBlockProps ============================== */

type FileBlockProps struct {
	BaseProps
	Url          string `json:"url" validate:"omitempty,isurl"`
	Caption      string `json:"caption,omitempty" validate:"omitempty,isfileblockcaption"`
	Name         string `json:"name,omitempty" validate:"omitempty,isfileblockname"`
	Size         int64  `json:"size,omitempty" validate:"omitempty,min=0"`
	PreviewWidth int    `json:"previewWidth,omitempty" validate:"omitempty"`
}

func (fbp *FileBlockProps) IsBlockProps() bool { return true }

func (fbp *FileBlockProps) Validate() error { return validation.Validator.Struct(fbp) }

/* ============================== ImageBlock ============================== */

type ImageBlockProps struct {
	FileBlockProps
}

func (ibp *ImageBlockProps) IsBlockProps() bool { return true }

func (ibp *ImageBlockProps) Validate() error { return validation.Validator.Struct(ibp) }

/* ============================== VideoBlock ============================== */

type VideoBlockProps struct {
	FileBlockProps
}

func (vbp *VideoBlockProps) IsBlockProps() bool { return true }

func (vbp *VideoBlockProps) Validate() error { return validation.Validator.Struct(vbp) }

/* ============================== AudioBlock ============================== */

type AudioBlockProps struct {
	FileBlockProps
}

func (abp *AudioBlockProps) IsBlockProps() bool { return true }

func (abp *AudioBlockProps) Validate() error { return validation.Validator.Struct(abp) }

/* ============================== CodeBlockProps ============================== */

type CodeBlockProps struct {
	BaseProps
	Language string `json:"language,omitempty" validate:"omitempty,isprogramminglanguage"`
}

func (cbp *CodeBlockProps) IsBlockProps() bool { return true }

func (cbp *CodeBlockProps) Validate() error { return validation.Validator.Struct(cbp) }

/* ============================== TableCellProps ============================== */

type TableCellProps struct {
	BaseProps
	RowSpan int `json:"rowspan" validate:"omitempty"`
	ColSpan int `json:"colspan" validate:"omitempty"`
}

func (tcp *TableCellProps) IsBlockProps() bool { return true }

func (tcp *TableCellProps) Validate() error { return validation.Validator.Struct(tcp) }

/* ============================== TableProps ============================== */

type TableProps struct {
	BaseProps
}

func (tp *TableProps) IsBlockProps() bool { return true }

func (tp *TableProps) Validate() error { return validation.Validator.Struct(tp) }

func ParseProps(blockType string, rawJSON []byte) (BlockProps, error) {
	if len(rawJSON) == 0 || string(rawJSON) == "null" {
		rawJSON = []byte("{}")
	}

	var props BlockProps

	switch blockType {
	case "heading":
		props = &HeadingProps{}
	case "checkListItem":
		props = &CheckListItemProps{}
	case "file":
		props = &FileBlockProps{}
	case "image":
		props = &ImageBlockProps{}
	case "video":
		props = &VideoBlockProps{}
	case "audio":
		props = &AudioBlockProps{}
	case "codeBlock":
		props = &CodeBlockProps{}
	case "table":
		props = &TableProps{}
	case "tableCell":
		props = &TableCellProps{}
	case "paragraph", "bulletListItem", "numberedListItem":
		props = &BaseProps{}
	default:
		props = &BaseProps{}
	}

	if err := json.Unmarshal(rawJSON, props); err != nil {
		return nil, err
	}

	if err := props.Validate(); err != nil {
		return nil, err
	}

	return props, nil
}
