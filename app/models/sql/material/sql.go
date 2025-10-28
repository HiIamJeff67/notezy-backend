package materialsql

import (
	_ "embed"
)

//go:embed get_my_material_and_its_parent_by_id.sql
var GetMyMaterialAndItsParentByIdSQL string
