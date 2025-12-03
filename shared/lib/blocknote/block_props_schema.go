package blocknote

type BaseProps struct {
	TextColor       string `json:"textColor,omitempty" validate:"omitempty,ishexcode"`
	BackgroundColor string `json:"backgroundColor,omitempty" validate:"omitempty,ishexcode"`
	TextAlignment   string `json:"textAlignment,omitempty" validate:"omitempty,istextalignment"`
}

type HeadingProps struct {
	BaseProps
	Level        int  `json:"level" validate:"required,isheadinglevel"`
	IsToggleable bool `json:"isToggleable,omitempty" validate:"omitempty"`
}

type CheckListItemProps struct {
	BaseProps
	Checked bool `json:"checked,omitempty"`
}

type FileBlockProps struct {
	BaseProps
	Url          string `json:"url" validate:"required,isurl"`
	Caption      string `json:"caption,omitempty" validate:"omitempty,isfileblockcaption"`
	Name         string `json:"name,omitempty" validate:"omitempty,isfileblockname"`
	Size         int64  `json:"size,omitempty" validate:"omitempty,min=0"`
	PreviewWidth int    `json:"previewWidth,omitempty" validate:"omitempty"`
}

type CodeBlockProps struct {
	BaseProps
	Language string `json:"language,omitempty" validate:"omitempty,isprogramminglanguage"`
}

type TableProps struct {
	BaseProps
}
