package postgres

type RelationName string

func (relationName RelationName) String() string {
	return string(relationName)
}
