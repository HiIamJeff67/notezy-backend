package types

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
	blocknote "github.com/HiIamJeff67/notezy-backend/shared/lib/blocknote"
)

type ArborizedEditableBlock struct {
	Id       uuid.UUID                `json:"id" validate:"required"`
	Type     enums.BlockType          `json:"type" validate:"required,isblocktype"`
	Props    blocknote.BlockProps     `json:"-"`
	Content  blocknote.BlockContent   `json:"-"`
	Children []ArborizedEditableBlock `json:"children" validate:"omitempty"`
}

func (block *ArborizedEditableBlock) UnmarshalJSON(data []byte) error {
	type alias ArborizedEditableBlock

	input := &struct {
		Props   json.RawMessage `json:"props"`
		Content json.RawMessage `json:"content"`
		*alias
	}{
		alias: (*alias)(block),
	}
	if err := json.Unmarshal(data, input); err != nil {
		return err
	}

	props, err := blocknote.ParseProps(string(block.Type), input.Props)
	if err != nil {
		return err
	}
	block.Props = props

	content := bytes.TrimSpace(input.Content)
	if len(content) == 0 || string(content) == "null" {
		return nil
	}

	switch content[0] {
	case '[':
		var inlineContent blocknote.InlineContentList
		if err := json.Unmarshal(content, &inlineContent); err != nil {
			return err
		}
		block.Content = inlineContent
	case '{':
		var tableContent blocknote.TableContent
		if err := json.Unmarshal(content, &tableContent); err != nil {
			return err
		}
		block.Content = &tableContent
	default:
		return errors.New("invalid content format: must be array or object")
	}

	return nil
}

func (block ArborizedEditableBlock) MarshalJSON() ([]byte, error) {
	type alias ArborizedEditableBlock

	return json.Marshal(&struct {
		Props   blocknote.BlockProps   `json:"props"`
		Content blocknote.BlockContent `json:"content"`
		*alias
	}{
		Props:   block.Props,
		Content: block.Content,
		alias:   (*alias)(&block),
	})
}

type RawFlattenedEditableBlock struct {
	Id            uuid.UUID       `json:"id"`
	ParentBlockId *uuid.UUID      `json:"parentBlockId"`
	PrevBlockId   *uuid.UUID      `json:"prevBlockId"`
	NextBlockId   *uuid.UUID      `json:"nextBlockId"`
	Type          enums.BlockType `json:"type"`
	Props         datatypes.JSON  `json:"props"`
	Content       datatypes.JSON  `json:"content"`
}
