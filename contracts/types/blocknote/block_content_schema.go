package blocknote

// BlockContent = []InlineContent | TableContent | PlainContent | undefined
type BlockContent interface {
	IsBlockContent() bool
}

type PlainContent string

func (pc PlainContent) IsBlockContent() bool { return true }
