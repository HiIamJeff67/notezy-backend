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

/* ============================== Arborized Editable Block ============================== */

// ArborizedEditableBlock is a type for frontend convience, it allowed the frontend to directly put the block output from the block note editor
// to this data struct, at the backend we can also simply unmarshal and validate the data struct
//
// To use it, you have to create a dto, and instead of embedding the ArborizedEditableBlock to the dto, we need to put the ArborizedEditableBlock as a type of a field in the dto
//
//	ex.
//	type CreateBlockReqDto {
//		ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock"`
//	    BlockGroupId uuid.UUID `json:"blockGroupId"`
//		ParentBlockId *uuid.UUID `json:"parentBlockId"`
//	}
type ArborizedEditableBlock struct {
	Id       uuid.UUID                `json:"id" validate:"required"`
	Type     enums.BlockType          `json:"type" validate:"required,isblocktype"`
	Props    blocknote.BlockProps     `json:"-"`
	Content  blocknote.BlockContent   `json:"-"`
	Children []ArborizedEditableBlock `json:"children" validate:"omitempty"`
}

func (aeb *ArborizedEditableBlock) UnmarshalJSON(data []byte) error {
	type AliasBlockDto ArborizedEditableBlock
	aux := &struct {
		Props   json.RawMessage `json:"props"`   // unmarshal to json raw message later temporarily
		Content json.RawMessage `json:"content"` // unmarshal to json raw message later temporarily
		*AliasBlockDto
	}{
		AliasBlockDto: (*AliasBlockDto)(aeb),
	}

	if err := json.Unmarshal(data, &aux); err != nil { // get the type in the Alias type of block dto
		return err
	}

	props, err := blocknote.ParseProps(string(aeb.Type), []byte("{}"))
	if err != nil {
		return err
	}
	aeb.Props = props

	trimContent := bytes.TrimSpace(aux.Content)

	if len(trimContent) > 0 && string(trimContent) != "null" {
		switch trimContent[0] {
		case '[':
			var list blocknote.InlineContentList
			if err := json.Unmarshal(trimContent, &list); err != nil {
				return err
			}
			// we have called the Validate() in the UnmarshalJSON() of InlineContentList for validating while unmarshaling the recursive data structure
			aeb.Content = list

		case '{':
			var table blocknote.TableContent
			if err := json.Unmarshal(trimContent, &table); err != nil {
				return err
			}
			// we have called the Validate() in the UnmarshalJSON() of TableContent for validating while unmarshaling the recursive data structure
			aeb.Content = &table

		default:
			return errors.New("invalid content format: must be array or object")
		}
	}

	return nil
}

func (aeb ArborizedEditableBlock) MarshalJSON() ([]byte, error) {
	type Alias ArborizedEditableBlock
	return json.Marshal(&struct {
		Props   blocknote.BlockProps   `json:"props"`
		Content blocknote.BlockContent `json:"content"`
		*Alias
	}{
		Props:   aeb.Props,
		Content: aeb.Content,
		Alias:   (*Alias)(&aeb),
	})
}

// RawArborizedEditableBlock is a type that used ONLY as a sub dto of the response dto,
// because we have make sure all the data of block props or content is type safe
// and valid before storing to the database, so we should trust the data coming
// from the database without any reason, and then just focus on the type safe
// and data validation while creating or updating the block props or content.
// Hence, the RawArborizedEditableBlock doesn't do any validation on any fields.
type RawArborizedEditableBlock struct {
	Id       uuid.UUID                   `json:"id"`
	Type     enums.BlockType             `json:"type"`
	Props    datatypes.JSON              `json:"props"`
	Content  datatypes.JSON              `json:"content"`
	Children []RawArborizedEditableBlock `json:"children"`
}

/* ============================== Flattened Editable Block ============================== */

type FlattenedEditableBlock struct {
	Id            uuid.UUID              `json:"id" validate:"required"`
	ParentBlockId *uuid.UUID             `json:"parentBlockId" validate:"omitempty"`
	Type          enums.BlockType        `json:"type" validate:"required,isblocktype"`
	Props         blocknote.BlockProps   `json:"-"`
	Content       blocknote.BlockContent `json:"-"`
}

func (feb *FlattenedEditableBlock) UnmarshalJSON(data []byte) error {
	type AliasBlockDto FlattenedEditableBlock
	aux := &struct {
		Props   json.RawMessage `json:"props"`   // unmarshal to json raw message later temporarily
		Content json.RawMessage `json:"content"` // unmarshal to json raw message later temporarily
		*AliasBlockDto
	}{
		AliasBlockDto: (*AliasBlockDto)(feb),
	}

	if err := json.Unmarshal(data, &aux); err != nil { // get the type in the Alias type of block dto
		return err
	}

	props, err := blocknote.ParseProps(string(feb.Type), []byte("{}"))
	if err != nil {
		return err
	}
	feb.Props = props

	trimContent := bytes.TrimSpace(aux.Content)

	if len(trimContent) > 0 && string(trimContent) != "null" {
		switch trimContent[0] {
		case '[':
			var list blocknote.InlineContentList
			if err := json.Unmarshal(trimContent, &list); err != nil {
				return err
			}
			// we have called the Validate() in the UnmarshalJSON() of InlineContentList for validating while unmarshaling the recursive data structure
			feb.Content = list

		case '{':
			var table blocknote.TableContent
			if err := json.Unmarshal(trimContent, &table); err != nil {
				return err
			}
			// we have called the Validate() in the UnmarshalJSON() of TableContent for validating while unmarshaling the recursive data structure
			feb.Content = &table

		default:
			return errors.New("invalid content format: must be array or object")
		}
	}

	return nil
}

func (feb FlattenedEditableBlock) MarshalJSON() ([]byte, error) {
	type Alias FlattenedEditableBlock
	return json.Marshal(&struct {
		Props   blocknote.BlockProps   `json:"props"`
		Content blocknote.BlockContent `json:"content"`
		*Alias
	}{
		Props:   feb.Props,
		Content: feb.Content,
		Alias:   (*Alias)(&feb),
	})
}

type RawFlattenedEditableBlock struct {
	Id            uuid.UUID       `json:"id"`
	ParentBlockId *uuid.UUID      `json:"parentBlockId"`
	Type          enums.BlockType `json:"type"`
	Props         datatypes.JSON  `json:"props"`
	Content       datatypes.JSON  `json:"content"`
}
