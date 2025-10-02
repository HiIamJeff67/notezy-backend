package shelfmaterialcascadingtriggersql

import (
	_ "embed"
)

//go:embed cascading_soft_delete_root_shelf_trigger.sql
var CascadingSoftDeleteRootShelfTriggerSQL string

//go:embed cascading_soft_delete_sub_shelf_trigger.sql
var CascadingSoftDeleteSubShelfTriggerSQL string

//go:embed cascading_restore_soft_deleted_root_shelf_trigger.sql
var CascadingRestoreRootShelfTriggerSQL string

//go:embed cascading_restore_soft_deleted_sub_shelf_trigger.sql
var CascadingRestoreSubShelfTriggerSQL string

//go:embed cascading_move_sub_shelf_trigger.sql
var CascadingMoveSubShelfTriggerSQL string
