package views

import _ "embed"

var (
	//go:embed user_view.sql
	userViewSQL string

	MigratingViewSQLs = []string{userViewSQL}
)
