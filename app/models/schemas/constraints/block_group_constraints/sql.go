package blockgroupconstraints

import (
	_ "embed"
)

var (
	//go:embed block_pack_id_prev_block_group_id_idx.sql
	BlockPackIdPrevBlockGroupIdIndexSQL string
)
