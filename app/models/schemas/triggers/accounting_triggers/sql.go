package accountingtriggersql

import (
	_ "embed"
)

var (
	//go:embed accounting_mutated_block_pack_trigger.sql
	AccountingMutatedBlockPackTriggerSQL string

	//go:embed accounting_mutated_block_trigger.sql
	AccountingMutatedBlockTriggerSQL string

	//go:embed accounting_mutated_root_shelf_trigger.sql
	AccountingMutatedRootShelfTriggerSQL string

	//go:embed accounting_mutated_sub_shelf_trigger.sql
	AccountingMutatedSubShelfTriggerSQL string

	//go:embed accouting_mutated_material_trigger.sql
	AccountingMutatedMaterialTriggerSQL string
)
