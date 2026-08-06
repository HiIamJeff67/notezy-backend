package database

type TableName string

func (tableName TableName) String() string {
	return string(tableName)
}
