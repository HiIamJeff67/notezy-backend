package adapters

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	blocknote "github.com/HiIamJeff67/notezy-backend/shared/lib/blocknote"
	queue "github.com/HiIamJeff67/notezy-backend/shared/lib/queue"
)

type EditableBlockAdapterInterface interface {
	Flatten(root *dtos.ArborizedEditableBlock) ([]dtos.FlattenedEditableBlock, *exceptions.Exception)
	FlattenToRaw(root *dtos.ArborizedEditableBlock) ([]dtos.RawFlattenedEditableBlock, int64, *exceptions.Exception)
	FlattenManyToRaw(roots []dtos.ArborizedEditableBlock) ([]dtos.RawFlattenedEditableBlock, int64, *exceptions.Exception)
	FlattenRawToRaw(root *dtos.RawArborizedEditableBlock) ([]dtos.RawFlattenedEditableBlock, int64, *exceptions.Exception)
	Arborize(root *dtos.FlattenedEditableBlock, childrenMap map[uuid.UUID][]dtos.FlattenedEditableBlock) (*dtos.ArborizedEditableBlock, *exceptions.Exception)
	ArborizeRawToRaw(root *dtos.RawFlattenedEditableBlock, childrenMap map[uuid.UUID][]dtos.RawFlattenedEditableBlock) (*dtos.RawArborizedEditableBlock, int64, *exceptions.Exception)
}

type EditableBlockAdapter struct{}

func NewEditableBlockAdapter() EditableBlockAdapterInterface {
	return &EditableBlockAdapter{}
}

// method to convert(flatten) arborized editable block to flattened editable block
func (ebca *EditableBlockAdapter) Flatten(
	root *dtos.ArborizedEditableBlock,
) ([]dtos.FlattenedEditableBlock, *exceptions.Exception) {
	if root == nil {
		return []dtos.FlattenedEditableBlock{}, nil
	}

	resultBlocks := []dtos.FlattenedEditableBlock{}
	resultBlocks = append(resultBlocks, dtos.FlattenedEditableBlock{
		Id:            root.Id,
		ParentBlockId: nil,
		Type:          root.Type,
		Props:         root.Props,
		Content:       root.Content,
	})
	q := queue.NewQueue[*dtos.ArborizedEditableBlock](constants.MAX_INT)
	q.Enqueue(root)
	visited := make(map[uuid.UUID]bool)
	for !q.IsEmpty() {
		current, err := q.Dequeue()
		if err != nil {
			return nil, exceptions.DataStructureLib.FailedToManipulateQueue().WithOrigin(err)
		}

		if visited[current.Id] {
			continue
		}
		visited[current.Id] = true

		for index, child := range current.Children {
			var prevBlockId *uuid.UUID
			if index > 0 {
				prev := current.Children[index-1].Id
				prevBlockId = &prev
			}

			var nextBlockId *uuid.UUID
			if index+1 < len(current.Children) {
				next := current.Children[index+1].Id
				nextBlockId = &next
			}

			resultBlocks = append(resultBlocks, dtos.FlattenedEditableBlock{
				Id:            child.Id,
				ParentBlockId: &current.Id,
				PrevBlockId:   prevBlockId,
				NextBlockId:   nextBlockId,
				Type:          child.Type,
				Props:         child.Props,
				Content:       child.Content,
			})

			q.Enqueue(&child)
		}
	}

	return resultBlocks, nil
}

// method to convert(flatten) arborized editable block to raw flattened editable block
func (ebca *EditableBlockAdapter) FlattenToRaw(
	root *dtos.ArborizedEditableBlock,
) ([]dtos.RawFlattenedEditableBlock, int64, *exceptions.Exception) {
	if root == nil {
		return []dtos.RawFlattenedEditableBlock{}, 0, nil
	}

	rootProps, err := json.Marshal(root.Props)
	if err != nil {
		return nil, 0, exceptions.Block.InvalidDto().WithOrigin(err)
	}
	rootContent, err := json.Marshal(root.Content)
	if err != nil {
		return nil, 0, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	var totalSize int64 = int64(len(rootProps) + len(rootContent))
	resultBlocks := []dtos.RawFlattenedEditableBlock{}
	resultBlocks = append(resultBlocks, dtos.RawFlattenedEditableBlock{
		Id:            root.Id,
		ParentBlockId: nil,
		Type:          root.Type,
		Props:         datatypes.JSON(rootProps),
		Content:       datatypes.JSON(rootContent),
	})
	q := queue.NewQueue[*dtos.ArborizedEditableBlock](constants.MAX_INT)
	q.Enqueue(root)
	visited := make(map[uuid.UUID]bool)
	for !q.IsEmpty() {
		current, err := q.Dequeue()
		if err != nil {
			return nil, 0, exceptions.DataStructureLib.FailedToManipulateQueue().WithOrigin(err)
		}

		if visited[current.Id] {
			continue
		}
		visited[current.Id] = true

		for index, child := range current.Children {
			props, err := json.Marshal(child.Props)
			if err != nil {
				return nil, 0, exceptions.Block.InvalidDto().WithOrigin(err)
			}
			content, err := json.Marshal(child.Content)
			if err != nil {
				return nil, 0, exceptions.Block.InvalidDto().WithOrigin(err)
			}

			totalSize += int64(len(props) + len(content))

			var prevBlockId *uuid.UUID
			if index > 0 {
				prev := current.Children[index-1].Id
				prevBlockId = &prev
			}

			var nextBlockId *uuid.UUID
			if index+1 < len(current.Children) {
				next := current.Children[index+1].Id
				nextBlockId = &next
			}

			resultBlocks = append(resultBlocks, dtos.RawFlattenedEditableBlock{
				Id:            child.Id,
				ParentBlockId: &current.Id,
				PrevBlockId:   prevBlockId,
				NextBlockId:   nextBlockId,
				Type:          child.Type,
				Props:         datatypes.JSON(props),
				Content:       datatypes.JSON(content),
			})

			q.Enqueue(&child)
		}
	}

	return resultBlocks, totalSize, nil
}

func (ebca *EditableBlockAdapter) FlattenManyToRaw(
	roots []dtos.ArborizedEditableBlock,
) ([]dtos.RawFlattenedEditableBlock, int64, *exceptions.Exception) {
	if len(roots) == 0 {
		return []dtos.RawFlattenedEditableBlock{}, 0, nil
	}

	type flattenItem struct {
		Block         *dtos.ArborizedEditableBlock
		ParentBlockId *uuid.UUID
		PrevBlockId   *uuid.UUID
		NextBlockId   *uuid.UUID
	}

	resultBlocks := make([]dtos.RawFlattenedEditableBlock, 0, len(roots))
	q := queue.NewQueue[flattenItem](constants.MAX_INT)
	for index := range roots {
		var prevBlockId *uuid.UUID
		if index > 0 {
			prev := roots[index-1].Id
			prevBlockId = &prev
		}

		var nextBlockId *uuid.UUID
		if index+1 < len(roots) {
			next := roots[index+1].Id
			nextBlockId = &next
		}

		q.Enqueue(flattenItem{
			Block:       &roots[index],
			PrevBlockId: prevBlockId,
			NextBlockId: nextBlockId,
		})
	}

	visited := make(map[uuid.UUID]bool)
	var totalSize int64
	for !q.IsEmpty() {
		current, err := q.Dequeue()
		if err != nil {
			return nil, 0, exceptions.DataStructureLib.FailedToManipulateQueue().WithOrigin(err)
		}
		if current.Block == nil || current.Block.Id == uuid.Nil || visited[current.Block.Id] {
			return nil, 0, exceptions.Block.InvalidDto()
		}
		visited[current.Block.Id] = true

		props, err := json.Marshal(current.Block.Props)
		if err != nil {
			return nil, 0, exceptions.Block.InvalidDto().WithOrigin(err)
		}
		content, err := json.Marshal(current.Block.Content)
		if err != nil {
			return nil, 0, exceptions.Block.InvalidDto().WithOrigin(err)
		}
		totalSize += int64(len(props) + len(content))

		resultBlocks = append(resultBlocks, dtos.RawFlattenedEditableBlock{
			Id:            current.Block.Id,
			ParentBlockId: current.ParentBlockId,
			PrevBlockId:   current.PrevBlockId,
			NextBlockId:   current.NextBlockId,
			Type:          current.Block.Type,
			Props:         datatypes.JSON(props),
			Content:       datatypes.JSON(content),
		})

		for index := range current.Block.Children {
			var prevBlockId *uuid.UUID
			if index > 0 {
				prev := current.Block.Children[index-1].Id
				prevBlockId = &prev
			}

			var nextBlockId *uuid.UUID
			if index+1 < len(current.Block.Children) {
				next := current.Block.Children[index+1].Id
				nextBlockId = &next
			}

			parentBlockId := current.Block.Id
			q.Enqueue(flattenItem{
				Block:         &current.Block.Children[index],
				ParentBlockId: &parentBlockId,
				PrevBlockId:   prevBlockId,
				NextBlockId:   nextBlockId,
			})
		}
	}

	return resultBlocks, totalSize, nil
}

// method to convert(flatten) raw arborized editable block to raw flattened editable block
func (ebca *EditableBlockAdapter) FlattenRawToRaw(
	root *dtos.RawArborizedEditableBlock,
) ([]dtos.RawFlattenedEditableBlock, int64, *exceptions.Exception) {
	if root == nil {
		return []dtos.RawFlattenedEditableBlock{}, 0, nil
	}

	var totalSize int64 = int64(len(root.Props) + len(root.Content))
	resultBlocks := []dtos.RawFlattenedEditableBlock{}
	resultBlocks = append(resultBlocks, dtos.RawFlattenedEditableBlock{
		Id:            root.Id,
		ParentBlockId: nil,
		Type:          root.Type,
		Props:         root.Props,
		Content:       root.Content,
	})
	q := queue.NewQueue[*dtos.RawArborizedEditableBlock](constants.MAX_INT)
	q.Enqueue(root)
	visited := make(map[uuid.UUID]bool)
	for !q.IsEmpty() {
		current, err := q.Dequeue()
		if err != nil {
			return nil, 0, exceptions.DataStructureLib.FailedToManipulateQueue().WithOrigin(err)
		}

		if visited[current.Id] {
			continue
		}
		visited[current.Id] = true

		for index, child := range current.Children {
			totalSize += int64(len(child.Props) + len(child.Content))

			var prevBlockId *uuid.UUID
			if index > 0 {
				prev := current.Children[index-1].Id
				prevBlockId = &prev
			}

			var nextBlockId *uuid.UUID
			if index+1 < len(current.Children) {
				next := current.Children[index+1].Id
				nextBlockId = &next
			}

			resultBlocks = append(resultBlocks, dtos.RawFlattenedEditableBlock{
				Id:            child.Id,
				ParentBlockId: &current.Id,
				PrevBlockId:   prevBlockId,
				NextBlockId:   nextBlockId,
				Type:          child.Type,
				Props:         child.Props,
				Content:       child.Content,
			})

			q.Enqueue(&child)
		}
	}

	return resultBlocks, totalSize, nil
}

// method to convert(arborize) flattened editable block to arborized editable block
func (ebca *EditableBlockAdapter) Arborize(
	root *dtos.FlattenedEditableBlock,
	childrenMap map[uuid.UUID][]dtos.FlattenedEditableBlock,
) (*dtos.ArborizedEditableBlock, *exceptions.Exception) {
	if root == nil {
		return &dtos.ArborizedEditableBlock{}, nil
	}

	rootContent := root.Content
	if root.Type == enums.BlockType_Table {
		if bytes, err := json.Marshal(root.Content); err == nil {
			var tableContent blocknote.TableContent
			if err := json.Unmarshal(bytes, &tableContent); err == nil {
				rootContent = &tableContent
			}
		}
	}
	resultRoot := dtos.ArborizedEditableBlock{
		Id:       root.Id,
		Type:     root.Type,
		Props:    root.Props,
		Content:  rootContent,
		Children: []dtos.ArborizedEditableBlock{},
	}
	q := queue.NewQueue[*dtos.ArborizedEditableBlock](len(childrenMap))
	q.Enqueue(&resultRoot)
	visited := make(map[uuid.UUID]bool, len(childrenMap))
	for !q.IsEmpty() {
		current, err := q.Dequeue()
		if err != nil {
			return nil, exceptions.DataStructureLib.FailedToManipulateQueue().WithOrigin(err)
		}

		if visited[current.Id] {
			continue
		}
		visited[current.Id] = true

		children, exist := childrenMap[current.Id]
		if !exist {
			continue
		}
		current.Children = make([]dtos.ArborizedEditableBlock, 0, len(children))
		for _, child := range children {
			arborizedEditableBlock := dtos.ArborizedEditableBlock{
				Id:       child.Id,
				Type:     child.Type,
				Props:    child.Props,
				Children: []dtos.ArborizedEditableBlock{}, // the children of the child should be initialize here
			}

			arborizedEditableBlock.Content = child.Content
			if child.Type == enums.BlockType_Table {
				if bytes, err := json.Marshal(child.Content); err == nil {
					var tableContent blocknote.TableContent
					if err := json.Unmarshal(bytes, &tableContent); err == nil {
						arborizedEditableBlock.Content = &tableContent
					}
				}
			}

			current.Children = append(current.Children, arborizedEditableBlock)
			currentChildPtr := &current.Children[len(current.Children)-1] // get the pointer to the child in current.Children
			q.Enqueue(currentChildPtr)
		}
	}

	return &resultRoot, nil
}

// method to convert(flatten) raw flattened editable block to arborized editable block
func (ebca *EditableBlockAdapter) ArborizeRawToRaw(
	root *dtos.RawFlattenedEditableBlock,
	childrenMap map[uuid.UUID][]dtos.RawFlattenedEditableBlock,
) (*dtos.RawArborizedEditableBlock, int64, *exceptions.Exception) {
	if root == nil {
		return &dtos.RawArborizedEditableBlock{}, 0, nil
	}

	var rootContent datatypes.JSON = root.Content
	if root.Type == enums.BlockType_Table {
		var tableContent blocknote.TableContent
		if err := json.Unmarshal(root.Content, &tableContent); err == nil {
			if bytes, err := json.Marshal(tableContent); err == nil {
				rootContent = datatypes.JSON(bytes)
			}
		}
	}
	var totalSize int64 = int64(len(root.Props) + len(rootContent))
	resultRoot := dtos.RawArborizedEditableBlock{
		Id:       root.Id,
		Type:     root.Type,
		Props:    root.Props,
		Content:  rootContent,
		Children: []dtos.RawArborizedEditableBlock{},
	}
	q := queue.NewQueue[*dtos.RawArborizedEditableBlock](len(childrenMap))
	q.Enqueue(&resultRoot)
	visited := make(map[uuid.UUID]bool, len(childrenMap))
	for !q.IsEmpty() {
		current, err := q.Dequeue()
		if err != nil {
			return nil, 0, exceptions.DataStructureLib.FailedToManipulateQueue().WithOrigin(err)
		}
		// at this point, current cannot be nil, because the root is not nil, and the below new element enqueued to the queue is alos not nil

		if visited[current.Id] {
			continue
		}
		visited[current.Id] = true

		children, exist := childrenMap[current.Id]
		if !exist {
			// no children under the current
			continue
		}
		current.Children = make([]dtos.RawArborizedEditableBlock, 0, len(children))
		for _, child := range children {
			arborizedEditableBlock := dtos.RawArborizedEditableBlock{
				Id:       child.Id,
				Type:     child.Type,
				Props:    child.Props,
				Children: []dtos.RawArborizedEditableBlock{}, // the children of the child should be initialize here
			}

			arborizedEditableBlock.Content = child.Content
			if child.Type == enums.BlockType_Table {
				var tableContent blocknote.TableContent
				if err := json.Unmarshal(child.Content, &tableContent); err == nil {
					if bytes, err := json.Marshal(tableContent); err == nil {
						arborizedEditableBlock.Content = datatypes.JSON(bytes)
					}
				}
			}

			totalSize += int64(len(arborizedEditableBlock.Props) + len(arborizedEditableBlock.Content))
			current.Children = append(current.Children, arborizedEditableBlock)
			currentChildPtr := &current.Children[len(current.Children)-1] // get the pointer to the child in current.Children
			q.Enqueue(currentChildPtr)                                    // make sure we passing the pointer of the editable child to the queue, so that we can modify its children field later
		}
	}

	return &resultRoot, totalSize, nil
}
