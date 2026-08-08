package apicontract

const (
	GetMyRootShelfByIdOperation           = "root-shelf.get-my-root-shelf-by-id"
	CreateRootShelfOperation              = "root-shelf.create"
	CreateRootShelvesOperation            = "root-shelf.create-many"
	UpdateMyRootShelfByIdOperation        = "root-shelf.update"
	UpdateMyRootShelvesByIdsOperation     = "root-shelf.update-many"
	RestoreMyRootShelfByIdOperation       = "root-shelf.restore"
	RestoreMyRootShelvesByIdsOperation    = "root-shelf.restore-many"
	DeleteMyRootShelfByIdOperation        = "root-shelf.delete"
	DeleteMyRootShelvesByIdsOperation     = "root-shelf.delete-many"
	GetMyRootShelfPermissionOperation     = "root-shelf.permission.get"
	CreateMyRootShelfPermissionOperation  = "root-shelf.permission.create"
	UpsertMyRootShelfPermissionOperation  = "root-shelf.permission.upsert"
	UpsertMyRootShelfPermissionsOperation = "root-shelf.permission.upsert-many"
	UpdateMyRootShelfPermissionOperation  = "root-shelf.permission.update"
	TransferMyRootShelfOwnershipOperation = "root-shelf.ownership.transfer"
	DeleteMyRootShelfPermissionOperation  = "root-shelf.permission.delete"
	DeleteMyRootShelfPermissionsOperation = "root-shelf.permission.delete-many"
	LeaveMyRootShelfOperation             = "root-shelf.membership.leave"
	LeaveMyRootShelvesOperation           = "root-shelf.membership.leave-many"
	SearchRootShelvesOperation            = "graphql.search-root-shelves"
)
