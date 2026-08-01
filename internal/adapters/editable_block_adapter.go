package adapters

import (
	"net/http"

	"gorm.io/datatypes"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/blocks"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas/enums"
	editableblock "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/editableblock"
)

type EditableBlockAdapterInterface interface {
	FlattenToRaw(root *blocksdto.ArborizedEditableBlock) ([]blocksdto.RawFlattenedEditableBlock, int64, *exceptions.Exception)
	FlattenManyToRaw(roots []blocksdto.ArborizedEditableBlock) ([]blocksdto.RawFlattenedEditableBlock, int64, *exceptions.Exception)
}

type EditableBlockAdapter struct{}

func NewEditableBlockAdapter() EditableBlockAdapterInterface {
	return &EditableBlockAdapter{}
}

/* ============================== Flatten Methods ============================== */

func (a *EditableBlockAdapter) FlattenToRaw(
	root *blocksdto.ArborizedEditableBlock,
) ([]blocksdto.RawFlattenedEditableBlock, int64, *exceptions.Exception) {
	if root == nil {
		return []blocksdto.RawFlattenedEditableBlock{}, 0, nil
	}

	portableRoot := toPortableEditableBlock(*root)
	blocks, totalSize, err := editableblock.FlattenToRaw(&portableRoot)
	if err != nil {
		return nil, 0, exceptions.New(
			"InvalidEditableBlock",
			"Block",
			"FlattenToRaw",
			"The editable block is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}

	return toRawFlattenedEditableBlocks(blocks), totalSize, nil
}

func (a *EditableBlockAdapter) FlattenManyToRaw(
	roots []blocksdto.ArborizedEditableBlock,
) ([]blocksdto.RawFlattenedEditableBlock, int64, *exceptions.Exception) {
	if len(roots) == 0 {
		return []blocksdto.RawFlattenedEditableBlock{}, 0, nil
	}

	portableRoots := make([]editableblock.ArborizedBlock, len(roots))
	for index := range roots {
		portableRoots[index] = toPortableEditableBlock(roots[index])
	}

	blocks, totalSize, err := editableblock.FlattenManyToRaw(portableRoots)
	if err != nil {
		return nil, 0, exceptions.New(
			"InvalidEditableBlock",
			"Block",
			"FlattenManyToRaw",
			"The editable blocks are invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}

	return toRawFlattenedEditableBlocks(blocks), totalSize, nil
}

/* ============================== Auxiliary Functions ============================== */

func toPortableEditableBlock(block blocksdto.ArborizedEditableBlock) editableblock.ArborizedBlock {
	children := make([]editableblock.ArborizedBlock, len(block.Children))
	for index := range block.Children {
		children[index] = toPortableEditableBlock(block.Children[index])
	}

	return editableblock.ArborizedBlock{
		Id:       block.Id,
		Type:     string(block.Type),
		Props:    block.Props,
		Content:  block.Content,
		Children: children,
	}
}

func toRawFlattenedEditableBlocks(blocks []editableblock.RawFlattenedBlock) []blocksdto.RawFlattenedEditableBlock {
	flattenedBlocks := make([]blocksdto.RawFlattenedEditableBlock, len(blocks))
	for index := range blocks {
		flattenedBlocks[index] = blocksdto.RawFlattenedEditableBlock{
			Id:            blocks[index].Id,
			ParentBlockId: blocks[index].ParentBlockId,
			PrevBlockId:   blocks[index].PrevBlockId,
			NextBlockId:   blocks[index].NextBlockId,
			Type:          enums.BlockType(blocks[index].Type),
			Props:         datatypes.JSON(blocks[index].Props),
			Content:       datatypes.JSON(blocks[index].Content),
		}
	}

	return flattenedBlocks
}
