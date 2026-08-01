package managementsql

import (
	_ "embed"
)

var (
	//go:embed get_all_enums.sql
	GetAllEnumsSQL string
)
