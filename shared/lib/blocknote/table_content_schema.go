package blocknote

import (
	"encoding/json"
	"fmt"

	validation "notezy-backend/app/validation"
)

/* ============================== TableCell & TableRow ============================== */

type TableCell InlineContentList

type TableRow struct {
	Cells []TableCell `json:"cells" validate:"required,min=1,max=100"`
}

/* ============================== TableContent ============================== */

type TableContentType string

const TableContentType_TableContent = "tableContent"

type TableContent struct {
	Type TableContentType `json:"type" validate:"required,eq=tableContent"`
	Rows []TableRow       `json:"rows" validate:"required,min=1,max=200"`
}

func (tc *TableContent) IsBlockContent() bool { return true }

func (tc *TableContent) Validate() error {
	if err := validation.Validator.Struct(tc); err != nil {
		return err
	}

	// instead of validating the rows and cells recursively,
	// we directly iterate the entire table by rows and cells then calling the Validate() method of the InlineContent inside of each cell

	expectedCells := -1

	for rowIndex, row := range tc.Rows {
		currentCells := len(row.Cells)
		if currentCells == 0 {
			return fmt.Errorf("table row %d is empty", rowIndex)
		}
		if expectedCells == -1 {
			expectedCells = currentCells
		} else if expectedCells != currentCells {
			return fmt.Errorf("jagged table detected: row %d has %d cells, expected %d", rowIndex, currentCells, expectedCells)
		}
	}

	return nil
}

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

	if err := tc.Validate(); err != nil {
		return err
	}

	return nil
}

func (tc *TableContent) MarshalJSON() ([]byte, error) {
	type AliaxTableContent TableContent
	var aux AliaxTableContent
	return json.Marshal(aux)
}
