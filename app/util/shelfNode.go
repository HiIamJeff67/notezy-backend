package util

import (
	"notezy-backend/app/exceptions"
	"regexp"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
)

type ShelfNode struct {
	Id          string                `json:"id" validate:"required"`
	Name        string                `json:"name" validate:"required,min=1,max=64"`
	Parent      *ShelfNode            `json:"parent" validate:"omitempty"`
	Children    map[string]*ShelfNode `json:"children" validate:"omitempty"`
	MaterialIds map[uuid.UUID]bool    `json:"materialIds" validate:"omitempty"` // leaves
}

/* ============================== Constructor ============================== */

func NewShelfNode(name string, parent *ShelfNode) *ShelfNode {
	if !isValidShelfName(name) {

	}
	shelfNodeId := GenerateSnowflakeID()
	return &ShelfNode{
		Id:          shelfNodeId,
		Name:        name,
		Parent:      parent,
		Children:    make(map[string]*ShelfNode),
		MaterialIds: make(map[uuid.UUID]bool),
	}
}

/* ============================== Private Methods ============================== */

func isValidShelfName(name string) bool {
	illegal := regexp.MustCompile(`[\/\\:\*\?"<>\|]`)
	return !illegal.MatchString(name)
}

func hasCycle(cur *ShelfNode, visited map[string]bool) bool {
	if cur == nil {
		return false
	}

	if visited[cur.Id] {
		return true
	}
	visited[cur.Id] = true

	for _, child := range cur.Children {
		if hasCycle(child, visited) {
			return true
		}
	}

	visited[cur.Id] = false
	return false
}

// Check if the `target` is a child of `cur`
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

// Check if the `target` is a parent of `cur`
func isParent(cur *ShelfNode, target *ShelfNode) bool {
	if cur == nil {
		return false
	}
	if cur == target {
		return true
	}
	return isParent(cur.Parent, target)
}

/* ============================== Public Methods ============================== */

func EncodeShelfNode(node *ShelfNode) ([]byte, *exceptions.Exception) {
	result, err := msgpack.Marshal(node)
	if err != nil {
		return nil, exceptions.Shelf.FailedToEncode(node).WithError(err)
	}
	return result, nil
}

func DecodeShelfNode(data []byte) (*ShelfNode, *exceptions.Exception) {
	var node ShelfNode
	if err := msgpack.Unmarshal(data, &node); err != nil {
		return nil, exceptions.Shelf.FailedToDecode(data).WithError(err)
	}
	return &node, nil
}

func (node *ShelfNode) GetAllParents() []*ShelfNode {
	if node == nil {
		return make([]*ShelfNode, 0)
	}

	var parents []*ShelfNode
	var cur *ShelfNode = node
	for cur != nil {
		cur = cur.Parent
		parents = append(parents, cur)
	}
	return parents
}

func (node *ShelfNode) GetAllParentIds() []string {
	if node == nil {
		return make([]string, 0)
	}

	var parents []string
	var cur *ShelfNode = node
	for cur != nil {
		cur = cur.Parent
		parents = append(parents, cur.Id)
	}
	return parents
}

func (node *ShelfNode) GetAllParentIdsInSet() map[string]bool {
	if node == nil {
		return make(map[string]bool)
	}

	parentsSet := make(map[string]bool)
	var cur *ShelfNode = node
	cur = cur.Parent
	for cur != nil {
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

	visited := map[string]bool{}
	return hasCycle(node, visited)
}

// Check if the current node as `node` has a child of the given node as `target`.
// If `node` does have a child of `target`, it will return true else false.
func (node *ShelfNode) HasChildOf(target *ShelfNode) bool {
	if node == nil {
		return false
	}

	return isParent(target, node)
}

// Check if current node as `node` is a child of the the given node as `target`.
// If `node` is a child of `target`, it will return true else false.
func (node *ShelfNode) IsChildOf(target *ShelfNode) bool {
	if node == nil {
		return false
	}

	return isParent(node, target)
}

// Check if the current node as `node` contains a subpath as `path`.
// If `path` is a subpath started from the `root`, it will return true else false.
// Note that the path only contains the id of the ShelfNode.
func (root *ShelfNode) IsSubpath(path []string) bool {
	if root == nil {
		return false
	}

	var cur *ShelfNode = &ShelfNode{
		Children: map[string]*ShelfNode{
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

func (destination *ShelfNode) InsertShelfNodes(targets []*ShelfNode) *exceptions.Exception {
	if destination == nil {
		return exceptions.Shelf.CallingMethodsWithNilValue()
	}

	parentsSet := destination.GetAllParentIdsInSet()
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

func (destination *ShelfNode) InsertMaterialById(targetId uuid.UUID) *exceptions.Exception {
	if destination == nil {
		return exceptions.Shelf.CallingMethodsWithNilValue()
	}

	destination.MaterialIds[targetId] = true
	return nil
}

func (destination *ShelfNode) InsertMaterialsByIds(targetIds []uuid.UUID) *exceptions.Exception {
	if destination == nil {
		return exceptions.Shelf.CallingMethodsWithNilValue()
	}

	for _, targetId := range targetIds {
		destination.MaterialIds[targetId] = true
	}
	return nil
}
