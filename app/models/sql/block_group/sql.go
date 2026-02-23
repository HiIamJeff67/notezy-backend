package blockgroupsql

import (
	_ "embed"
)

var (
	//go:embed get_my_block_group_and_its_blocks_by_id.sql
	GetMyBlockGroupAndItsBlocksByIdSQL string

	//go:embed get_my_block_groups_and_their_blocks_by_block_pack_id.sql
	GetMyBlockGroupsAndTheirBlocksByBlockPackIdSQL string

	//go:embed collect_garbage_block_group_by_ids.sql
	CollectGarbageBlockGroupByIdsSQL string
)
