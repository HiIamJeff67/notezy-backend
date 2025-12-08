package adapters

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
	queue "notezy-backend/shared/lib/queue"
)

type EditableBlockAdapterInterface interface {
	Flatten(root *dtos.ArborizedEditableBlock) ([]dtos.FlattenedEditableBlock, *exceptions.Exception)
	FlattenToRaw(root *dtos.ArborizedEditableBlock) ([]dtos.RawFlattenedEditableBlock, *exceptions.Exception)
	FlattenRawToRaw(root *dtos.RawArborizedEditableBlock) ([]dtos.RawFlattenedEditableBlock, *exceptions.Exception)
	Arborize(root *dtos.FlattenedEditableBlock, childrenMap map[uuid.UUID][]dtos.FlattenedEditableBlock) (*dtos.ArborizedEditableBlock, *exceptions.Exception)
	ArborizeRawToRaw(root *dtos.RawFlattenedEditableBlock, childrenMap map[uuid.UUID][]dtos.RawFlattenedEditableBlock) (*dtos.RawArborizedEditableBlock, *exceptions.Exception)
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
		return nil, nil
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
			return nil, exceptions.DataStructureLib.FailedToManipulateQueue().WithError(err)
		}

		if visited[current.Id] {
			continue
		}
		visited[current.Id] = true

		for _, child := range current.Children {
			resultBlocks = append(resultBlocks, dtos.FlattenedEditableBlock{
				Id:            child.Id,
				ParentBlockId: &current.Id,
				Type:          child.Type,
				Props:         child.Props,
				Content:       child.Content,
			})

			q.Enqueue(&child)
		}
	}

	return resultBlocks, nil
}

// method to conver(flatten) arborized editable block to raw flattened editable block
func (ebca *EditableBlockAdapter) FlattenToRaw(
	root *dtos.ArborizedEditableBlock,
) ([]dtos.RawFlattenedEditableBlock, *exceptions.Exception) {
	if root == nil {
		return nil, nil
	}

	rootProps, err := json.Marshal(root.Props)
	if err != nil {
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}
	rootContent, err := json.Marshal(root.Content)
	if err != nil {
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}

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
			return nil, exceptions.DataStructureLib.FailedToManipulateQueue().WithError(err)
		}

		if visited[current.Id] {
			continue
		}
		visited[current.Id] = true

		for _, child := range current.Children {
			props, err := json.Marshal(child.Props)
			if err != nil {
				return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
			}
			content, err := json.Marshal(child.Content)
			if err != nil {
				return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
			}

			resultBlocks = append(resultBlocks, dtos.RawFlattenedEditableBlock{
				Id:            root.Id,
				ParentBlockId: nil,
				Type:          root.Type,
				Props:         datatypes.JSON(props),
				Content:       datatypes.JSON(content),
			})

			q.Enqueue(&child)
		}
	}

	return resultBlocks, nil
}

// method to convert(flatten) raw arborized editable block to raw flattened editable block
func (ebca *EditableBlockAdapter) FlattenRawToRaw(
	root *dtos.RawArborizedEditableBlock,
) ([]dtos.RawFlattenedEditableBlock, *exceptions.Exception) {
	if root == nil {
		return nil, nil
	}

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
			return nil, exceptions.DataStructureLib.FailedToManipulateQueue().WithError(err)
		}

		if visited[current.Id] {
			continue
		}
		visited[current.Id] = true

		for _, child := range current.Children {
			resultBlocks = append(resultBlocks, dtos.RawFlattenedEditableBlock{
				Id:            child.Id,
				ParentBlockId: &current.Id,
				Type:          child.Type,
				Props:         child.Props,
				Content:       child.Content,
			})

			q.Enqueue(&child)
		}
	}

	return resultBlocks, nil
}

// method to convert(arborize) flattened editable block to arborized editable block
func (ebca *EditableBlockAdapter) Arborize(
	root *dtos.FlattenedEditableBlock,
	childrenMap map[uuid.UUID][]dtos.FlattenedEditableBlock,
) (*dtos.ArborizedEditableBlock, *exceptions.Exception) {
	if root == nil {
		return nil, nil
	}

	resultRoot := dtos.ArborizedEditableBlock{
		Id:       root.Id,
		Type:     root.Type,
		Props:    root.Props,
		Content:  root.Content,
		Children: []dtos.ArborizedEditableBlock{},
	}
	q := queue.NewQueue[*dtos.ArborizedEditableBlock](len(childrenMap))
	q.Enqueue(&resultRoot)
	visited := make(map[uuid.UUID]bool, len(childrenMap))
	for !q.IsEmpty() {
		current, err := q.Dequeue()
		if err != nil {
			return nil, exceptions.DataStructureLib.FailedToManipulateQueue().WithError(err)
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
			current.Children = append(current.Children, dtos.ArborizedEditableBlock{
				Id:       child.Id,
				Type:     child.Type,
				Props:    child.Props,
				Content:  child.Content,
				Children: []dtos.ArborizedEditableBlock{}, // the children of the child should be initialize here
			})
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
) (*dtos.RawArborizedEditableBlock, *exceptions.Exception) {
	if root == nil {
		return nil, nil
	}

	resultRoot := dtos.RawArborizedEditableBlock{
		Id:       root.Id,
		Type:     root.Type,
		Props:    root.Props,
		Content:  root.Content,
		Children: []dtos.RawArborizedEditableBlock{},
	}
	q := queue.NewQueue[*dtos.RawArborizedEditableBlock](len(childrenMap))
	q.Enqueue(&resultRoot)
	visited := make(map[uuid.UUID]bool, len(childrenMap))
	for !q.IsEmpty() {
		current, err := q.Dequeue()
		if err != nil {
			return nil, exceptions.DataStructureLib.FailedToManipulateQueue().WithError(err)
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
			current.Children = append(current.Children, dtos.RawArborizedEditableBlock{
				Id:       child.Id,
				Type:     child.Type,
				Props:    child.Props,
				Content:  child.Content,
				Children: []dtos.RawArborizedEditableBlock{}, // the children of the child should be initialize here
			})
			currentChildPtr := &current.Children[len(current.Children)-1] // get the pointer to the child in current.Children
			q.Enqueue(currentChildPtr)                                    // make sure we passing the pointer of the editable child to the queue, so that we can modify its children field later
		}
	}

	return &resultRoot, nil
}
