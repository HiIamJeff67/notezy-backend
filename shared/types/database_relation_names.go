package types

// Use this type to define a specific type of a database tables' relations
type RelationName string

func (tr RelationName) String() string {
	return string(tr)
}
