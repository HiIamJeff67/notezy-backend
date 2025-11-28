package blocknote

type BaseProps struct {
	TextColor       string `json:"textColor,omitempty"`
	BackgroundColor string `json:"backgroundColor,omitempty"`
	TextAlignment   string `json:"textAlignment,omitempty"` // "left", "center", "right", "justify"
}

type HeadingProps struct {
	BaseProps
	Level        int  `json:"level" validate:"required,min=1,max=3"`
	IsToggleable bool `json:"isToggleable" validate:"omitempty"`
}

type CheckListItemProps struct {
	BaseProps
	Checked bool `json:"checked,omitempty"`
}

type FileBlockProps struct {
	BaseProps
	Url          string `json:"url" validate:"required"`
	Caption      string `json:"caption,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int64  `json:"size,omitempty"`
	PreviewWidth int    `json:"previewWidth,omitempty"`
}

type CodeBlockProps struct {
	BaseProps
	Language string `json:"language,omitempty"` // "javascript", "go", etc.
}

type TableProps struct {
	BaseProps
}
