package editableblock

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type ArborizedBlock struct {
	Id       uuid.UUID
	Type     string
	Props    any
	Content  any
	Children []ArborizedBlock
}

type RawFlattenedBlock struct {
	Id            uuid.UUID
	ParentBlockId *uuid.UUID
	PrevBlockId   *uuid.UUID
	NextBlockId   *uuid.UUID
	Type          string
	Props         json.RawMessage
	Content       json.RawMessage
}

type flattenItem struct {
	block         *ArborizedBlock
	parentBlockId *uuid.UUID
	prevBlockId   *uuid.UUID
	nextBlockId   *uuid.UUID
}

func FlattenToRaw(root *ArborizedBlock) ([]RawFlattenedBlock, int64, error) {
	if root == nil {
		return []RawFlattenedBlock{}, 0, nil
	}

	return FlattenManyToRaw([]ArborizedBlock{*root})
}

func FlattenManyToRaw(roots []ArborizedBlock) ([]RawFlattenedBlock, int64, error) {
	if len(roots) == 0 {
		return []RawFlattenedBlock{}, 0, nil
	}

	queue := make([]flattenItem, 0, len(roots))
	for index := range roots {
		var prevBlockId *uuid.UUID
		if index > 0 {
			previousId := roots[index-1].Id
			prevBlockId = &previousId
		}

		var nextBlockId *uuid.UUID
		if index+1 < len(roots) {
			nextId := roots[index+1].Id
			nextBlockId = &nextId
		}

		queue = append(queue, flattenItem{
			block:       &roots[index],
			prevBlockId: prevBlockId,
			nextBlockId: nextBlockId,
		})
	}

	flattenedBlocks := make([]RawFlattenedBlock, 0, len(roots))
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

		flattenedBlocks = append(flattenedBlocks, RawFlattenedBlock{
			Id:            item.block.Id,
			ParentBlockId: item.parentBlockId,
			PrevBlockId:   item.prevBlockId,
			NextBlockId:   item.nextBlockId,
			Type:          item.block.Type,
			Props:         props,
			Content:       content,
		})

		for index := range item.block.Children {
			var prevBlockId *uuid.UUID
			if index > 0 {
				previousId := item.block.Children[index-1].Id
				prevBlockId = &previousId
			}

			var nextBlockId *uuid.UUID
			if index+1 < len(item.block.Children) {
				nextId := item.block.Children[index+1].Id
				nextBlockId = &nextId
			}

			parentBlockId := item.block.Id
			queue = append(queue, flattenItem{
				block:         &item.block.Children[index],
				parentBlockId: &parentBlockId,
				prevBlockId:   prevBlockId,
				nextBlockId:   nextBlockId,
			})
		}
	}

	return flattenedBlocks, totalSize, nil
}
