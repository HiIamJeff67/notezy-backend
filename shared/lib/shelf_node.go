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
	Parent      *ShelfNode               `json:"parent"` // if the ShelfNode is a root, the Parent field MUST be nil
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
	parent *ShelfNode,
) (*ShelfNode, *exceptions.Exception) {
	if !IsValidShelfName(name) {
		return nil, exceptions.Shelf.FailedToConstructNewShelfNode("name")
	}

	shelfNodeId := uuid.New()
	result := &ShelfNode{
		Id:          shelfNodeId,
		Name:        name,
		Parent:      parent,
		Children:    make(map[uuid.UUID]*ShelfNode),
		MaterialIds: make(map[uuid.UUID]bool),
	}
	return result, nil
}

// Note: The below time and space analysis will use N as the number of ShelfNode, and M as the number of material ids in the entire tree
/* ============================== Private Methods ============================== */

// Check if the `target` is a child of `cur`.
func isChild(cur *ShelfNode, target *ShelfNode) bool {
	if cur == nil {
		return false
	}
	if cur == target {
		return true
	}
	for _, child := range cur.Children {
		if isChild(child, target) {
			return true
		}
	}
	return false
}

// Check if the `target` is a parent of `cur`.
func isParent(root *ShelfNode, target *ShelfNode, ctx context.Context) bool {
	if root == nil || target == nil {
		return false
	}

	cur := root
	traverseCount := 0

	for cur != nil {
		if traverseCount%constants.CheckPointPerTraverse == 0 {
			select {
			case <-ctx.Done():
				return true
			default:
				// no-op
			}
		}
		traverseCount++

		if cur == target {
			return true
		}

		cur = cur.Parent
	}
	return false
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
	if node.Parent != nil {
		return nil, exceptions.Shelf.CannotEncodeNonRootShelfNode(node)
	}

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

// Get all the parent of the current node and returning the instance as a list.
func (node *ShelfNode) GetAllParents(ctx context.Context) []*ShelfNode {
	if node == nil {
		return make([]*ShelfNode, 0)
	}

	var parents []*ShelfNode
	var cur *ShelfNode = node
	traverseCount := 0

	for cur != nil {
		if traverseCount%constants.CheckPointPerTraverse == 0 {
			select {
			case <-ctx.Done():
				return nil
			default:
				// no-op
			}
		}
		traverseCount++

		cur = cur.Parent
		parents = append(parents, cur)
	}
	return parents
}

// Get all the parent of the current node and returning their ids as a list.
func (node *ShelfNode) GetAllParentIds(ctx context.Context) []uuid.UUID {
	if node == nil {
		return make([]uuid.UUID, 0)
	}

	var parents []uuid.UUID
	var cur *ShelfNode = node
	traverseCount := 0

	for cur != nil {
		if traverseCount%constants.CheckPointPerTraverse == 0 {
			select {
			case <-ctx.Done():
				return nil
			default:
				// no-op
			}
		}

		cur = cur.Parent
		parents = append(parents, cur.Id)
	}
	return parents
}

// Get all the parent of the current node and returning their ids as a set.
func (node *ShelfNode) GetAllParentIdsInSet(ctx context.Context) map[uuid.UUID]bool {
	if node == nil {
		return make(map[uuid.UUID]bool)
	}

	parentsSet := make(map[uuid.UUID]bool)
	var cur *ShelfNode = node
	cur = cur.Parent
	traverseCount := 0

	for cur != nil {
		if traverseCount%constants.CheckPointPerTraverse == 0 {
			select {
			case <-ctx.Done():
				return nil
			default:
				// no-op
			}
		}

		parentsSet[cur.Id] = true
		cur = cur.Parent
	}
	return parentsSet
}

// Check if the children of the given ShelfNode have any cycles.
// If there's any cycles detected, it will return true else false.
func (node *ShelfNode) IsChildrenCircular() bool {
	if node == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.MaxTraverseTimeout)
	defer cancel()

	cur := node
	queue := make([]ShelfNodeWithDepth, 0)
	queue = append(queue, ShelfNodeWithDepth{Node: cur, Depth: 0})
	visited := map[uuid.UUID]bool{}
	traverseCount := 0

	for len(queue) > 0 {
		if traverseCount%constants.CheckPointPerTraverse == 0 {
			select {
			case <-ctx.Done():
				return true
			default:
				// no-op
			}
		}
		traverseCount++

		current := queue[0]
		queue = queue[1:] // pop

		if visited[current.Node.Id] {
			return true
		}
		visited[current.Node.Id] = true

		for _, child := range current.Node.Children {
			queue = append(queue, ShelfNodeWithDepth{
				Node:  child,
				Depth: current.Depth + 1,
			})
		}
	}

	return false
}

// Check if the current node as `node` has a child of the given node as `target`.
// If `node` does have a child of `target`, it will return true else false.
func (node *ShelfNode) HasChildOf(target *ShelfNode) bool {
	if node == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.MaxTraverseTimeout)
	defer cancel()

	return isParent(target, node, ctx)
}

// Check if current node as `node` is a child of the the given node as `target`.
// If `node` is a child of `target`, it will return true else false.
func (node *ShelfNode) IsChildOf(target *ShelfNode) bool {
	if node == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.MaxTraverseTimeout)
	defer cancel()

	return isParent(node, target, ctx)
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

// Traverse the entire shelf tree by using breadth first search and mean while analysis to collect some informations,
// return a generated summary including :
//   - total number of shelf nodes
//   - total number of materials
//   - max width of the children(ShelfNode) which is equal to m
//   - max depth of the children(ShelfNode)
//   - unique material ids in a list
//   - a exception if there's any error happened
func (node *ShelfNode) AnalysisAndGenerateSummary() (
	totalShelfNodes int,
	totalMaterials int,
	maxWidth int64,
	maxDepth int64,
	uniqueMaterialIds []uuid.UUID,
	exception *exceptions.Exception,
) {
	if node == nil {
		return 0, 0, 0, 0, make([]uuid.UUID, 0), exceptions.Shelf.CallingMethodsWithNilValue()
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.MaxTraverseTimeout)
	defer cancel()

	cur := node
	queue := make([]ShelfNodeWithDepth, 0)
	queue = append(queue, ShelfNodeWithDepth{Node: cur, Depth: 0})
	visited := map[uuid.UUID]bool{}
	uniqueMaterialIdsSet := map[uuid.UUID]bool{}
	// the return value of totalShelfNodes does the same thing as the traverseCount

	for len(queue) > 0 {
		if totalShelfNodes%constants.CheckPointPerTraverse == 0 {
			select {
			case <-ctx.Done():
				return totalShelfNodes,
					totalMaterials,
					maxWidth,
					maxDepth,
					uniqueMaterialIds,
					exceptions.Shelf.Timeout(constants.MaxTraverseTimeout)
			default:
				// no-op
			}
		}

		current := queue[0]
		queue = queue[1:] // pop
		totalShelfNodes++
		maxDepth = max(maxDepth, current.Depth)

		if visited[current.Node.Id] {
			return 0, 0, 0, 0, make([]uuid.UUID, 0), exceptions.Shelf.RepeatedShelfNodesDetected()
		}
		visited[current.Node.Id] = true

		for materialId, exist := range current.Node.MaterialIds {
			if !exist {
				continue
			}

			if uniqueMaterialIdsSet[materialId] {
				return 0, 0, 0, 0, make([]uuid.UUID, 0), exceptions.Shelf.RepeatedMaterialIdsDetected()
			}
			uniqueMaterialIdsSet[materialId] = true
			uniqueMaterialIds = append(uniqueMaterialIds, materialId)
		}

		maxWidth = max(maxWidth, int64(len(current.Node.Children)))
		for _, child := range current.Node.Children {
			queue = append(queue, ShelfNodeWithDepth{
				Node:  child,
				Depth: current.Depth + 1,
			})
		}
	}

	return totalShelfNodes,
		len(uniqueMaterialIds),
		maxWidth,
		maxDepth,
		uniqueMaterialIds,
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
	target.Parent = nil
	destination.Children[target.Id] = target
	target.Parent = destination
	return nil
}

// Insert the given list of ShelfNodes of `targets` into the `destination`,
// Note that all the ShelfNode of `targets` MUST NOT be one of the parents of `destination`
func (destination *ShelfNode) InsertShelfNodes(targets []*ShelfNode) *exceptions.Exception {
	if destination == nil {
		return exceptions.Shelf.CallingMethodsWithNilValue()
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.MaxTraverseTimeout)
	defer cancel()

	parentsSet := destination.GetAllParentIdsInSet(ctx)
	for _, target := range targets {
		if parentsSet[target.Id] {
			return exceptions.Shelf.InsertParentIntoItsChildren(destination, target)
		}
		target.Parent = nil
		destination.Children[target.Id] = target
		target.Parent = destination
	}
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
