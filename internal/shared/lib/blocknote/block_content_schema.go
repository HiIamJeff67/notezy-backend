package blocknote

// BlockContent = []InlineContent | TableContent | undefined
type BlockContent interface {
	IsBlockContent() bool
	Validate() error
}
