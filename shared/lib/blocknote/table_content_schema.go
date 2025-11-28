package blocknote

type TableContentType string

const TableContentType_TableContent = "tableContent"

type TableCell []InlineContent

type TableRow struct {
	Cells []TableCell `json:"cells"`
}

type TableContent struct {
	Type TableContentType `json:"type" validate:"required,eq=tableContent"`
	Rows []TableRow       `json:"rows" validate:"omitempty"`
}
