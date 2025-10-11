package blocknote

/* ==================== TableContent Definitions ==================== */

type TableCellType string

const TableCellType_TableCell TableCellType = "tableCell"

type TableCell struct {
	Type    TableCellType          `json:"type" validate:"required,eq=tableCell"`
	Props   map[string]interface{} `json:"props" validate:"omitempty"`
	Content []InlineContent        `json:"content" validate:"omitempty"`
}

type TableContentType string

const TableContentType_TableContent = "tableContent"

type TableContent struct {
	Type         TableContentType              `json:"type" validate:"required,eq=tableContent"`
	ColumnWidths []float32                     `json:"columnWidths" validate:"omitempty"`
	HeaderRows   float32                       `json:"headerRows" validate:"required"`
	Rows         []struct{ cells []TableCell } `json:"rows" validate:"omitempty"`
}
