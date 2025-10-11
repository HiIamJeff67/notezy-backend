package blocknote

import (
	"bytes"
	"encoding/json"
)

/* ==================== BlockContent ==================== */

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
		var arr []InlineContent
		if err := json.Unmarshal(trim, &arr); err != nil {
			return err
		}
		bc.InlineContent = arr
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
			var tbl TableContent
			if err := json.Unmarshal(trim, &tbl); err != nil {
				return err
			}
			bc.TableContent = &tbl
			bc.InlineContent = nil
			return nil
		}
		return json.Unmarshal(trim, &bc.InlineContent)
	default:
		return json.Unmarshal(trim, &bc.InlineContent)
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

/* ==================== Block Definitions ==================== */

type BlockType string

const BlockType_Paragraph = "paragraph"
const BlockType_Heading = "heading"
const BlockType_Quote = "quote"
const BlockType_BulletListItem = "bulletListItem"
const BlockType_NumberedListItem = "numberedListItem"
const BlockType_CheckListItem = "checkListItem"
const BlockType_ToggleListItemBlock = "toggleListItemBlock"
const BlockType_Table = "table"
const BlockType_File = "file"
const BlockType_Image = "image"
const BlockType_Video = "video"
const BlockType_Audio = "audio"
const BlockType_CodeBlock = "codeBlock"

type Block struct {
	Id       string                 `json:"id" validate:"required"`
	Type     BlockType              `json:"type" validate:"required"`
	Props    map[string]interface{} `json:"props" validate:"omitempty"`
	Content  *BlockContent          `json:"content" validate:"omitempty"`
	Children []Block                `json:"children" validate:"omitempty"`
}
