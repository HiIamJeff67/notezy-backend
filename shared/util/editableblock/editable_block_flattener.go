package editableblock

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	typescontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
)

type editableBlockFlattenItem struct {
	block         *typescontract.ArborizedEditableBlock
	parentBlockId *uuid.UUID
	prevBlockId   *uuid.UUID
	nextBlockId   *uuid.UUID
}

/* ============================== Auxiliary Functions ============================== */

func previousEditableBlockId(blocks []typescontract.ArborizedEditableBlock, index int) *uuid.UUID {
	if index == 0 {
		return nil
	}

	previousBlockId := blocks[index-1].Id
	return &previousBlockId
}

func nextEditableBlockId(blocks []typescontract.ArborizedEditableBlock, index int) *uuid.UUID {
	if index+1 == len(blocks) {
		return nil
	}

	nextBlockId := blocks[index+1].Id
	return &nextBlockId
}

/* ============================== Main Methods ============================== */

func FlattenEditableBlock(root *typescontract.ArborizedEditableBlock) ([]typescontract.RawFlattenedEditableBlock, int64, error) {
	if root == nil {
		return []typescontract.RawFlattenedEditableBlock{}, 0, nil
	}

	return FlattenEditableBlocks([]typescontract.ArborizedEditableBlock{*root})
}

func FlattenEditableBlocks(roots []typescontract.ArborizedEditableBlock) ([]typescontract.RawFlattenedEditableBlock, int64, error) {
	if len(roots) == 0 {
		return []typescontract.RawFlattenedEditableBlock{}, 0, nil
	}

	queue := make([]editableBlockFlattenItem, 0, len(roots))
	for index := range roots {
		queue = append(queue, editableBlockFlattenItem{
			block:       &roots[index],
			prevBlockId: previousEditableBlockId(roots, index),
			nextBlockId: nextEditableBlockId(roots, index),
		})
	}

	flattenedBlocks := make([]typescontract.RawFlattenedEditableBlock, 0, len(roots))
	visitedBlockIds := make(map[uuid.UUID]bool)
	var totalSize int64
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]
		if item.block == nil || item.block.Id == uuid.Nil {
			return nil, 0, fmt.Errorf("editable block id is required")
		}
		if visitedBlockIds[item.block.Id] {
			return nil, 0, fmt.Errorf("duplicate editable block id: %s", item.block.Id)
		}
		visitedBlockIds[item.block.Id] = true

		props, err := json.Marshal(item.block.Props)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal editable block props: %w", err)
		}
		content, err := json.Marshal(item.block.Content)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal editable block content: %w", err)
		}
		totalSize += int64(len(props) + len(content))

		flattenedBlocks = append(flattenedBlocks, typescontract.RawFlattenedEditableBlock{
			Id:            item.block.Id,
			ParentBlockId: item.parentBlockId,
			PrevBlockId:   item.prevBlockId,
			NextBlockId:   item.nextBlockId,
			Type:          item.block.Type,
			Props:         props,
			Content:       content,
		})

		for index := range item.block.Children {
			parentBlockId := item.block.Id
			queue = append(queue, editableBlockFlattenItem{
				block:         &item.block.Children[index],
				parentBlockId: &parentBlockId,
				prevBlockId:   previousEditableBlockId(item.block.Children, index),
				nextBlockId:   nextEditableBlockId(item.block.Children, index),
			})
		}
	}

	return flattenedBlocks, totalSize, nil
}
