package maintainingblocktriggers

import (
	_ "embed"
)

var (
	// go:embed maintaining_final_block_group_on_insert_trigger.sql
	MaintainingFinalBlockGroupOnInsertTriggerSQL string

	// go:embed maintaining_final_block_group_on_update_trigger.sql
	MaintainingFinalBlockGroupOnUpdateTriggerSQL string

	// go:embed maintaining_final_block_group_on_delete_trigger.sql
	MaintainingFinalBlockGroupOnDeleteTriggerSQL string
)
