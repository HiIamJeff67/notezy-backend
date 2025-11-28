package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

/* ============================== Block Type ============================== */

type BlockType string

const (
	BlockType_Paragraph BlockType = "paragraph"
	BlockType_Heading   BlockType = "heading"
	BlockType_Quote     BlockType = "quote"

	BlockType_BulletListItem   BlockType = "bulletListItem"
	BlockType_NumberedListItem BlockType = "numberedListItem"
	BlockType_CheckListItem    BlockType = "checkListItem"

	BlockType_Image BlockType = "image"
	BlockType_Video BlockType = "video"
	BlockType_Audio BlockType = "audio"
	BlockType_File  BlockType = "file"

	BlockType_Table     BlockType = "table"
	BlockType_CodeBlock BlockType = "codeBlock"
)

var AllBlockTypes = []BlockType{
	BlockType_Paragraph,
	BlockType_Heading,
	BlockType_BulletListItem,
	BlockType_NumberedListItem,
	BlockType_CheckListItem,
	BlockType_Image,
	BlockType_Video,
	BlockType_Audio,
	BlockType_File,
	BlockType_Table,
	BlockType_CodeBlock,
}

var AllBlockTypeStrings = []string{
	string(BlockType_Paragraph),
	string(BlockType_Heading),
	string(BlockType_BulletListItem),
	string(BlockType_NumberedListItem),
	string(BlockType_CheckListItem),
	string(BlockType_Image),
	string(BlockType_Video),
	string(BlockType_Audio),
	string(BlockType_File),
	string(BlockType_Table),
	string(BlockType_CodeBlock),
}

func (bt BlockType) Name() string {
	return reflect.TypeOf(bt).Name()
}

func (bt *BlockType) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*bt = BlockType(string(v))
		return nil
	case string:
		*bt = BlockType(v)
		return nil
	}
	return fmt.Errorf("cannot scan %T into BlockType", value)
}

func (bt BlockType) Value() (driver.Value, error) {
	return string(bt), nil
}

func (bt BlockType) String() string {
	return string(bt)
}

func (bt *BlockType) IsValidEnum() bool {
	return slices.Contains(AllBlockTypes, *bt)
}

func ConvertStringToBlockType(enumString string) (*BlockType, error) {
	for _, blockType := range AllBlockTypes {
		if string(blockType) == enumString {
			return &blockType, nil
		}
	}
	return nil, fmt.Errorf("invalid block type: %s", enumString)
}
