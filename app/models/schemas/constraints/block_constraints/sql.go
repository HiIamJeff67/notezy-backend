package blockconstraints

import (
	_ "embed"
)

var (
	//go:embed block_tree_root_idx.sql
	BlockTreeRootIndexSQL string
)
