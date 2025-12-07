package dtos

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "notezy-backend/app/models/schemas/enums"
	blocknote "notezy-backend/shared/lib/blocknote"
)

/* ============================== Auxiliary Data Form ============================== */

// EditableBlockContent is a type for frontend convience, it allowed the frontend to directly put the block output from the block note editor
// to this data struct, at the backend we can also simply unmarshal and validate the data struct
//
// To use it, you have to create a dto, and instead of embedding the EditableBlockContent to the dto, we need to put the EditableBlockContent as a type of a field in the dto
//
//	ex.
//	type CreateBlockReqDto {
//		EditableBlockContent EditableBlockContent `json:"editableBlockContent"`
//	    BlockGroupId uuid.UUID `json:"blockGroupId"`
//		ParentBlockId *uuid.UUID `json:"parentBlockId"`
//	}
type EditableBlockContent struct {
	Id       uuid.UUID              `json:"id" validate:"required"`
	Type     enums.BlockType        `json:"type" validate:"required"`
	Props    blocknote.BlockProps   `json:"-"`
	Content  blocknote.BlockContent `json:"-"`
	Children []EditableBlockContent `json:"children" validate:"omitempty"`
}

func (bd *EditableBlockContent) UnmarshalJSON(data []byte) error {
	type AliasBlockDto EditableBlockContent
	aux := &struct {
		Props   json.RawMessage `json:"props"`   // unmarshal to json raw message later temporarily
		Content json.RawMessage `json:"content"` // unmarshal to json raw message later temporarily
		*AliasBlockDto
	}{
		AliasBlockDto: (*AliasBlockDto)(bd),
	}

	if err := json.Unmarshal(data, &aux); err != nil { // get the type in the Alias type of block dto
		return err
	}

	props, err := blocknote.ParseProps(string(bd.Type), []byte("{}"))
	if err != nil {
		return err
	}
	bd.Props = props

	trimContent := bytes.TrimSpace(aux.Content)

	if len(trimContent) > 0 && string(trimContent) != "null" {
		switch trimContent[0] {
		case '[':
			var list blocknote.InlineContentList
			if err := json.Unmarshal(trimContent, &list); err != nil {
				return err
			}
			// we have called the Validate() in the UnmarshalJSON() of InlineContentList for validating while unmarshaling the recursive data structure
			bd.Content = list

		case '{':
			var table blocknote.TableContent
			if err := json.Unmarshal(trimContent, &table); err != nil {
				return err
			}
			// we have called the Validate() in the UnmarshalJSON() of TableContent for validating while unmarshaling the recursive data structure
			bd.Content = &table

		default:
			return errors.New("invalid content format: must be array or object")
		}
	}

	return nil
}

func (bd EditableBlockContent) MarshalJSON() ([]byte, error) {
	type Alias EditableBlockContent
	return json.Marshal(&struct {
		Props   blocknote.BlockProps   `json:"props"`
		Content blocknote.BlockContent `json:"content"`
		*Alias
	}{
		Props:   bd.Props,
		Content: bd.Content,
		Alias:   (*Alias)(&bd),
	})
}

// EditableRawBlockContent is a type that used ONLY as a sub dto of the response dto,
// because we have make sure all the data of block props or content is type safe
// and valid before storing to the database, so we should trust the data coming
// from the database without any reason, and then just focus on the type safe
// and data validation while creating or updating the block props or content
type EditableRawBlockContent struct {
	Id       uuid.UUID                 `json:"id"`
	Type     enums.BlockType           `json:"type"`
	Props    datatypes.JSON            `json:"props"`
	Content  datatypes.JSON            `json:"content"`
	Children []EditableRawBlockContent `json:"children"`
}

/* ============================== Request DTO ============================== */

/* ============================== Response DTO ============================== */
