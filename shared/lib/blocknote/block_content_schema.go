package blocknote

import (
	"bytes"
	"encoding/json"
	"errors"
)

// BlockContent = []InlineContent | TableContent | undefined
type BlockContent struct {
	InlineContent []InlineContent
	TableContent  *TableContent
}

func (bc *BlockContent) UnmarshalJSON(b []byte) error {
	trim := bytes.TrimSpace(b)
	if len(trim) == 0 || string(trim) == "null" {
		// treat it as undefined
		return nil
	}

	switch trim[0] {
	case '[': // detect []InlineContent, since it is the only type of BlockContent which is an array
		var content []InlineContent
		if err := json.Unmarshal(trim, &content); err != nil {
			return err
		}
		bc.InlineContent = content
		bc.TableContent = nil
		return nil
	case '{': // detect TableContent, since it is neither undefined nor []InlineContent
		var kind struct {
			Type TableContentType `json:"type"`
		}
		if err := json.Unmarshal(trim, &kind); err != nil {
			return err
		}
		if kind.Type == TableContentType_TableContent {
			var content TableContent
			if err := json.Unmarshal(trim, &content); err != nil {
				return err
			}
			bc.TableContent = &content
			bc.InlineContent = nil
			return nil
		}
		return errors.New("unknown block content object type")
	default:
		return errors.New("invalid json format for block content")
	}
}

func (bc BlockContent) MarshalJSON() ([]byte, error) {
	if bc.TableContent != nil {
		return json.Marshal(bc.TableContent)
	}
	if bc.InlineContent == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(bc.InlineContent)
}
