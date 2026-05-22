package materialconstraints

import (
	_ "embed"
)

var (
	//go:embed material_pure_file_migration.sql
	MaterialPureFileMigrationSQL string
)
