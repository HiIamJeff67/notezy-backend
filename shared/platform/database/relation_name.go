package database

type RelationName string

func (relationName RelationName) String() string {
	return string(relationName)
}
