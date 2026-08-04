package blocknote

import "encoding/json"

/* ============================== TableCell ============================== */

type TableCellType string

const TableCellType_TableCell TableCellType = "tableCell"

type TableCell struct {
	Type    TableCellType     `json:"type" validate:"required,eq=tableCell"`
	Content InlineContentList `json:"content" validate:"omitempty,dive"`
	Props   TableCellProps    `json:"props" validate:"omitempty"`
}

/* ============================== TableRow ============================== */

type TableRow struct {
	Cells []TableCell `json:"cells" validate:"required,min=1,max=100,dive"`
}

/* ============================== TableContent ============================== */

type TableContentType string

const TableContentType_TableContent = "tableContent"

type TableContent struct {
	Type         TableContentType `json:"type" validate:"required,eq=tableContent"`
	ColumnWidths []*string        `json:"columnWidths" validate:"omitempty"`
	Rows         []TableRow       `json:"rows" validate:"required,min=1,max=200,dive"`
}

func (tc *TableContent) IsBlockContent() bool { return true }

func (tc *TableContent) UnmarshalJSON(b []byte) error {
	// To avoid infinitly recursion happen on the UnmarshalJSON
	// 	ex. json.Unmarshal(tc) -> find tc has implement the Unmarshaler method of UnmarshalJSON -> calling UnmarshalJSON (the current method)
	//      -> inside the UnmarshalJSON method, we call the json.Unmarshal(tc) (instead of json.Unmarshal(aux)) -> find tc... etc => infinitly recursion
	// hence, we need the aliax type of table content
	type AliaxTableContent TableContent // the AliaxTableContent type does NOT have the implementation of the Unmarshaler (UnmarshalJSON() method)
	var aux AliaxTableContent
	if err := json.Unmarshal(b, &aux); err != nil { // remember we have implement the UnmarshalJSON() in the InlineContent, it'll be called to unmarshal and validate the inline content
		return err
	}

	*tc = TableContent(aux)

	return nil
}

func (tc *TableContent) MarshalJSON() ([]byte, error) {
	type AliaxTableContent TableContent
	aux := AliaxTableContent(*tc)
	return json.Marshal(aux)
}
