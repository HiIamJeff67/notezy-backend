package lib

import (
	"context"
	"regexp"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"

	exceptions "notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
)

// The node data type of the shelf which is present in M-way Tree structure.
type ShelfNode struct {
	Id          uuid.UUID                `json:"id"`
	Name        string                   `json:"name"`
	Children    map[uuid.UUID]*ShelfNode `json:"children"`
	MaterialIds map[uuid.UUID]bool       `json:"materialIds"` // leaves
}

type ShelfNodeWithDepth struct {
	Node  *ShelfNode
	Depth int64
}

/* ============================== Constructor ============================== */

func NewShelfNode(
	ownerId uuid.UUID, // since only the owner can create his/her own shelf
	name string,
) (*ShelfNode, *exceptions.Exception) {
	if !IsValidShelfName(name) {
		return nil, exceptions.Shelf.FailedToConstructNewShelfNode("name")
	}

	shelfNodeId := uuid.New()
	result := &ShelfNode{
		Id:          shelfNodeId,
		Name:        name,
		Children:    make(map[uuid.UUID]*ShelfNode),
		MaterialIds: make(map[uuid.UUID]bool),
	}
	return result, nil
}

// Note: The below time and space analysis will use N as the number of ShelfNode, and M as the number of material ids in the entire tree
/* ============================== Private Methods ============================== */

func (node *ShelfNode) estimateSingleNodeSize() int64 {
	// 36 bytes for UUID + 20 bytes for the approximate structure cost
	// + approximate estimate UTF-8 to be 2 bytes/char
	var size int64 = int64(56 + len(node.Name)*2)

	childrenCount := int64(len(node.Children))
	materialCount := int64(len(node.MaterialIds))
	size += (childrenCount + materialCount) * 36 // for each key of UUID data type
	size += (childrenCount + materialCount) * 1  // null/reference overhead and boolean value

	return size
}

// Check if the `target` is a child of `node`.
func isChild(node *ShelfNode, target *ShelfNode) (isChild bool, exception *exceptions.Exception) {
	if node == nil || target == nil {
		return false, exceptions.Shelf.CallingMethodsWithNilValue()
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.MaxShelfTreeTraverseTimeout)
	defer cancel()

	queue := make([]ShelfNodeWithDepth, 0)
	queue = append(queue, ShelfNodeWithDepth{Node: node, Depth: 0})
	visited := map[uuid.UUID]bool{}
	var traverseCount int32 = 0
	var maxWidth int32 = 0
	var maxDepth int32 = 0

	for len(queue) > 0 {
		levelSize := len(queue)
		maxWidth = max(maxWidth, int32(levelSize))
		maxDepth++

		if maxWidth > constants.MaxShelfTreeWidth {
			return false, exceptions.Shelf.MaximumWidthExceeded(maxWidth, constants.MaxShelfTreeWidth)
		}
		if maxDepth > constants.MaxShelfTreeDepth {
			return false, exceptions.Shelf.MaximumDepthExceeded(maxDepth, constants.MaxShelfTreeDepth)
		}
		for i := 0; i < levelSize; i++ {
			if traverseCount > constants.MaxNumOfShelfTreeTraversedNodes {
				return false, exceptions.Shelf.MaximumTraverseCountExceeded(traverseCount, constants.MaxNumOfShelfTreeTraversedNodes)
			}
			if traverseCount%constants.CheckPointPerShelfTreeTraverse == 0 {
				select {
				case <-ctx.Done():
					return false, exceptions.Shelf.Timeout(constants.MaxShelfTreeTraverseTimeout)
				default:
					// no-op
				}
			}
			traverseCount++

			current := queue[0]
			queue = queue[1:] // pop

			if visited[current.Node.Id] {
				return false, exceptions.Shelf.CircularChildrenDetectedInShelfNode()
			}
			visited[current.Node.Id] = true

			for _, child := range current.Node.Children {
				if child == target {
					return true, nil
				}
				queue = append(queue, ShelfNodeWithDepth{
					Node:  child,
					Depth: current.Depth + 1,
				})
			}
		}
	}

	return false, nil
}

/* ============================== Public Methods ============================== */

func IsValidShelfName(name string) bool {
	if len(name) > 128 {
		return false
	}
	return !regexp.MustCompile(`[\/\\:\*\?"<>\|]`).MatchString(name)
}

// Encode the entire shelf node INCLUDE its children but EXCLUDE its parent into a bytes type,
// This operation is done by using a msgpack 3rd party library.
func EncodeShelfNode(node *ShelfNode) ([]byte, *exceptions.Exception) {
	result, err := msgpack.Marshal(node)
	if err != nil {
		return nil, exceptions.Shelf.FailedToEncode(node).WithError(err)
	}
	return result, nil
}

// Decode the entire shelf node INCLUDE its children but EXCLUDE its parent into a bytes type,
// This operation is done by using a msgpack 3rd party library.
func DecodeShelfNode(data []byte) (*ShelfNode, *exceptions.Exception) {
	var node ShelfNode
	if err := msgpack.Unmarshal(data, &node); err != nil {
		return nil, exceptions.Shelf.FailedToDecode(data).WithError(err)
	}
	return &node, nil
}

// Check if the children of the given ShelfNode is a simple tree
// which means it shouldn't contain any cycles.
// If there's any cycles detected, it will return false else true.
// Note that if there's any other exceptions, it will also return false as well.
func (node *ShelfNode) IsChildrenSimple() (isSimple bool, exception *exceptions.Exception) {
	if node == nil {
		return false, exceptions.Shelf.CallingMethodsWithNilValue()
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.MaxShelfTreeTraverseTimeout)
	defer cancel()

	queue := make([]ShelfNodeWithDepth, 0)
	queue = append(queue, ShelfNodeWithDepth{Node: node, Depth: 0})
	visited := map[uuid.UUID]bool{}
	var traverseCount int32 = 0
	var maxWidth int32 = 0
	var maxDepth int32 = 0

	for len(queue) > 0 {
		levelSize := len(queue)
		maxWidth = max(maxWidth, int32(levelSize))
		maxDepth++

		if maxWidth > constants.MaxShelfTreeWidth {
			return false, exceptions.Shelf.MaximumWidthExceeded(maxWidth, constants.MaxShelfTreeWidth)
		}
		if maxDepth > constants.MaxShelfTreeDepth {
			return false, exceptions.Shelf.MaximumDepthExceeded(maxDepth, constants.MaxShelfTreeDepth)
		}

		for i := 0; i < levelSize; i++ {
			if traverseCount > constants.MaxNumOfShelfTreeTraversedNodes {
				return false, exceptions.Shelf.MaximumTraverseCountExceeded(traverseCount, constants.MaxNumOfShelfTreeTraversedNodes)
			}
			if traverseCount%constants.CheckPointPerShelfTreeTraverse == 0 {
				select {
				case <-ctx.Done():
					return false, exceptions.Shelf.Timeout(constants.MaxShelfTreeTraverseTimeout)
				default:
					// no-op
				}
			}

			current := queue[0]
			queue = queue[1:] // pop
			traverseCount++

			if visited[current.Node.Id] {
				return false, nil
			}
			visited[current.Node.Id] = true

			for _, child := range current.Node.Children {
				queue = append(queue, ShelfNodeWithDepth{
					Node:  child,
					Depth: current.Depth + 1,
				})
			}
		}

	}

	return true, nil
}

// Check if the current node as `node` has a child of the given node as `target`.
// If `node` does have a child of `target`, it will return true else false.
func (node *ShelfNode) HasChildOf(target *ShelfNode) bool {
	if node == nil {
		return false
	}

	isNodeAChildOfTarget, exception := isChild(node, target)
	if exception != nil {
		exception.Log()
		return false
	}

	return isNodeAChildOfTarget
}

// Check if current node as `node` is a child of the the given node as `target`.
// If `node` is a child of `target`, it will return true else false.
func (node *ShelfNode) IsChildOf(target *ShelfNode) bool {
	if node == nil {
		return false
	}

	isTargetAChildOfNode, exception := isChild(target, node)
	if exception != nil {
		exception.Log()
		return false
	}

	return isTargetAChildOfNode
}

// Check if the current node as `node` contains a subpath as `path`.
// If `path` is a subpath started from the `root`, it will return true else false.
// Note that the path only contains the id of the ShelfNode.
func (root *ShelfNode) HasSubpathOf(path []uuid.UUID) bool {
	if root == nil {
		return false
	}

	var cur *ShelfNode = &ShelfNode{
		Children: map[uuid.UUID]*ShelfNode{
			root.Id: root,
		},
	}

	for _, id := range path {
		child, ok := cur.Children[id]
		if !ok {
			return false
		}
		cur = child
	}

	return true
}

/* ============================== Methods for Services ============================== */

type ShelfSummary struct {
	EncodedStructureByteSize int64
	TotalShelfNodes          int32
	TotalMaterials           int32
	MaxWidth                 int32
	MaxDepth                 int32
	UniqueMaterialIds        []uuid.UUID
}

// Traverse the entire shelf tree by using breadth first search and mean while analysis to collect some informations,
// return a generated summary including :
//   - total number of shelf nodes
//   - total number of materials
//   - max width of the children(ShelfNode) which is equal to m
//   - max depth of the children(ShelfNode)
//   - unique material ids in a list
//   - a exception if there's any error happened
func (node *ShelfNode) AnalysisAndGenerateSummary() (
	summary *ShelfSummary,
	exception *exceptions.Exception,
) {
	if node == nil {
		return nil, exceptions.Shelf.CallingMethodsWithNilValue()
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.MaxShelfTreeTraverseTimeout)
	defer cancel()

	queue := make([]ShelfNodeWithDepth, 0)
	queue = append(queue, ShelfNodeWithDepth{Node: node, Depth: 0})
	visited := map[uuid.UUID]bool{}
	uniqueMaterialIdsSet := map[uuid.UUID]bool{}
	uniqueMaterialIds := make([]uuid.UUID, 0)
	var encodedStructureByteSize int64 = 0
	var totalShelfNodes int32 = 0
	var maxWidth int32 = 0
	var maxDepth int32 = 0
	// the return value of totalShelfNodes does the same thing as the traverseCount

	for len(queue) > 0 {
		levelSize := len(queue)
		maxWidth = max(maxWidth, int32(levelSize))
		maxDepth++

		if maxWidth > constants.MaxShelfTreeWidth {
			return nil, exceptions.Shelf.MaximumWidthExceeded(maxWidth, constants.MaxShelfTreeWidth)
		}
		if maxDepth > constants.MaxShelfTreeDepth {
			return nil, exceptions.Shelf.MaximumDepthExceeded(maxDepth, constants.MaxShelfTreeDepth)
		}

		for i := 0; i < levelSize; i++ {
			if totalShelfNodes > constants.MaxNumOfShelfTreeTraversedNodes {
				return nil, exceptions.Shelf.MaximumTraverseCountExceeded(totalShelfNodes, constants.MaxNumOfShelfTreeTraversedNodes)
			}
			if totalShelfNodes%constants.CheckPointPerShelfTreeTraverse == 0 {
				select {
				case <-ctx.Done():
					return nil, exceptions.Shelf.Timeout(constants.MaxShelfTreeTraverseTimeout)
				default:
					// no-op
				}
			}

			current := queue[0]
			queue = queue[1:] // pop
			totalShelfNodes++

			if visited[current.Node.Id] {
				return nil, exceptions.Shelf.RepeatedShelfNodesDetected()
			}
			visited[current.Node.Id] = true

			encodedStructureByteSize += current.Node.estimateSingleNodeSize()

			for materialId, exist := range current.Node.MaterialIds {
				if !exist {
					continue
				}

				if uniqueMaterialIdsSet[materialId] {
					return nil, exceptions.Shelf.RepeatedMaterialIdsDetected()
				}
				uniqueMaterialIdsSet[materialId] = true
				uniqueMaterialIds = append(uniqueMaterialIds, materialId)
			}

			for _, child := range current.Node.Children {
				queue = append(queue, ShelfNodeWithDepth{
					Node:  child,
					Depth: current.Depth + 1,
				})
			}
		}
	}

	return &ShelfSummary{
			encodedStructureByteSize,
			totalShelfNodes,
			int32(len(uniqueMaterialIds)),
			maxWidth,
			maxDepth,
			uniqueMaterialIds,
		},
		nil
}

// Insert the given ShelfNode of `target` into the `destination`,
// Note that `target` MUST NOT be one of the parents of `destination`
func (destination *ShelfNode) InsertShelfNode(target *ShelfNode) *exceptions.Exception {
	if destination == nil {
		return exceptions.Shelf.CallingMethodsWithNilValue()
	}

	if destination.IsChildOf(target) {
		return exceptions.Shelf.InsertParentIntoItsChildren(destination, target)
	}
	destination.Children[target.Id] = target
	return nil
}

// Insert the Material with id of `targetId` into the MaterialIds field of the `destination`,
// Note that this function MUST ONLY be called with passing the new Material
func (destination *ShelfNode) InsertMaterialById(targetId uuid.UUID) *exceptions.Exception {
	if destination == nil {
		return exceptions.Shelf.CallingMethodsWithNilValue()
	}

	destination.MaterialIds[targetId] = true
	return nil
}

// Insert the all the Materials with each of id inside the `targetIds` into the MaterialIds field of the `destination`,
// Note that this function MUST ONLY be called with passing the new Materials
func (destination *ShelfNode) InsertMaterialsByIds(targetIds []uuid.UUID) *exceptions.Exception {
	if destination == nil {
		return exceptions.Shelf.CallingMethodsWithNilValue()
	}

	for _, targetId := range targetIds {
		destination.MaterialIds[targetId] = true
	}
	return nil
}
