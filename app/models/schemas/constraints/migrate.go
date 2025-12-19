package constraints

import blockgroupconstraints "notezy-backend/app/models/schemas/constraints/block_group_constraints"

var MigratingConstraintSQLs = []string{
	blockgroupconstraints.BlockPackIdPrevBlockGroupIdIndexSQL,
}
