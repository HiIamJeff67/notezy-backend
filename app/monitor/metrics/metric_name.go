package metrics

/* ============================== requests fields ============================== */

type metricNameAuth struct {
	Register          string
	RegisterViaGoogle string
	Login             string
	LoginViaGoogle    string
	Logout            string
	SendAuthCode      string
	ValidateEmail     string
	ResetEmail        string
	ForgetPassword    string
	ResetMe           string
	DeleteMe          string
}

type metricNameUser struct {
	GetUserData string
	GetMe       string
	UpdateMe    string
}

type metricNameUserInfo struct {
	GetMyInfo    string
	UpdateMyInfo string
}

type metricNameUserSetting struct {
	GetMySetting string
}

type metricNameUserAccount struct {
	GetMyAccount        string
	UpdateMyAccount     string
	BindGoogleAccount   string
	UnbindGoogleAccount string
}

type metricNameRootShelf struct {
	GetMyRootShelfById        string
	SearchRecentRootShelves   string
	CreateRootShelf           string
	CreateRootShelves         string
	UpdateMyRootShelfById     string
	UpdateMyRootShelvesByIds  string
	RestoreMyRootShelfById    string
	RestoreMyRootShelvesByIds string
	DeleteMyRootShelfById     string
	DeleteMyRootShelvesByIds  string
}

type metricNameSubShelf struct {
	GetMySubShelfById                       string
	GetMySubShelvesByPrevSubShelfId         string
	GetAllMySubShelvesByRootShelfId         string
	GetMySubShelvesAndItemsByPrevSubShelfId string
	CreateSubShelfByRootShelfId             string
	CreateSubShelvesByRootShelfIds          string
	UpdateMySubShelfById                    string
	UpdateMySubShelvesByIds                 string
	MoveMySubShelf                          string
	MoveMySubShelves                        string
	BatchMoveMySubShelves                   string
	RestoreMySubShelfById                   string
	RestoreMySubShelvesByIds                string
	DeleteMySubShelfById                    string
	DeleteMySubShelvesByIds                 string
}

type metricNameMaterial struct {
	GetMyMaterialById                string
	GetMyMaterialAndItsParentById    string
	GetMyMaterialsByParentSubShelfId string
	GetAllMyMaterialsByRootShelfId   string
	CreateTextbookMaterial           string
	CreateNotebookMaterial           string
	UpdateMyMaterialById             string
	SaveMyNotebookMaterialById       string
	MoveMyMaterialById               string
	MoveMyMaterialsByIds             string
	RestoreMyMaterialById            string
	RestoreMyMaterialsByIds          string
	DeleteMyMaterialById             string
	DeleteMyMaterialsByIds           string
}

type metricNameBlockPack struct {
	GetMyBlockPackById                string
	GetMyBlockPackAndItsParentById    string
	GetMyBlockPacksByParentSubShelfId string
	GetAllMyBlockPacksByRootShelfId   string
	CreateBlockPack                   string
	CreateBlockPacks                  string
	UpdateMyBlockPackById             string
	UpdateMyBlockPacksByIds           string
	MoveMyBlockPackById               string
	MoveMyBlockPacksByIds             string
	BatchMoveMyBlockPacksByIds        string
	RestoreMyBlockPackById            string
	RestoreMyBlockPacksByIds          string
	DeleteMyBlockPackById             string
	DeleteMyBlockPacksByIds           string
}

type metricNameBlockGroup struct {
	GetMyBlockGroupById                                    string
	GetMyBlockGroupAndItsBlocksById                        string
	GetMyBlockGroupsAndTheirBlocksByIds                    string
	GetMyBlockGroupsAndTheirBlocksByBlockPackId            string
	GetMyBlockGroupsByPrevBlockGroupId                     string
	GetAllMyBlockGroupsByBlockPackId                       string
	InsertBlockGroupByBlockPackId                          string
	InsertBlockGroupAndItsBlocksByBlockPackId              string
	InsertBlockGroupsAndTheirBlocksByBlockPackId           string
	InsertSequentialBlockGroupsAndTheirBlocksByBlockPackId string
	MoveMyBlockGroupsByIds                                 string
	RestoreMyBlockGroupById                                string
	RestoreMyBlockGroupsByIds                              string
	DeleteMyBlockGroupById                                 string
	DeleteMyBlockGroupsByIds                               string
}

type metricNameBlock struct {
	GetMyBlockById             string
	GetMyBlocksByIds           string
	GetMyBlocksByBlockGroupId  string
	GetMyBlocksByBlockGroupIds string
	GetMyBlocksByBlockPackId   string
	GetAllMyBlocks             string
	InsertBlock                string
	InsertBlocks               string
	UpdateMyBlockById          string
	UpdateMyBlocksByIds        string
	RestoreMyBlockById         string
	RestoreMyBlocksByIds       string
	DeleteMyBlockById          string
	DeleteMyBlocksByIds        string
}

type MetricNameRequests struct {
	Total       string
	Auth        metricNameAuth
	User        metricNameUser
	UserInfo    metricNameUserInfo
	UserSetting metricNameUserSetting
	UserAccount metricNameUserAccount
	RootShelf   metricNameRootShelf
	SubShelf    metricNameSubShelf
	Material    metricNameMaterial
	BlockPack   metricNameBlockPack
	BlockGroup  metricNameBlockGroup
	Block       metricNameBlock
}

/* ============================== responses fields ============================== */

type MetricNameResponse struct {
	Success struct {
		Total string
	}
	Failed struct {
		Total        string
		Timeout      string
		Unauthorized string
		RateLimit    string
	}
	Email struct {
		Welcome  string
		AuthCode string
	}
}

var MetricNames = struct {
	Server struct {
		Requests  MetricNameRequests
		Responses MetricNameResponse
	}
}{
	Server: struct {
		Requests  MetricNameRequests
		Responses MetricNameResponse
	}{
		Requests: MetricNameRequests{
			Total: "server.requests.total",
			Auth: metricNameAuth{
				Register:          "server.requests.auth.register",
				RegisterViaGoogle: "server.requests.auth.registerViaGoogle",
				Login:             "server.requests.auth.login",
				LoginViaGoogle:    "server.requests.auth.loginViaGoogle",
				Logout:            "server.requests.auth.logout",
				SendAuthCode:      "server.requests.auth.sendAuthCode",
				ValidateEmail:     "server.requests.auth.validateEmail",
				ResetEmail:        "server.requests.auth.resetEmail",
				ForgetPassword:    "server.requests.auth.forgetPassword",
				ResetMe:           "server.requests.auth.resetMe",
				DeleteMe:          "server.requests.auth.deleteMe",
			},
			User: metricNameUser{
				GetUserData: "server.requests.user.getUserData",
				GetMe:       "server.requests.user.getMe",
				UpdateMe:    "server.requests.user.updateMe",
			},
			UserInfo: metricNameUserInfo{
				GetMyInfo:    "server.requests.userInfo.getMyInfo",
				UpdateMyInfo: "server.requests.userInfo.updateMyInfo",
			},
			UserSetting: metricNameUserSetting{
				GetMySetting: "server.requests.userSetting.getMySetting",
			},
			UserAccount: metricNameUserAccount{
				GetMyAccount:        "server.requests.userAccount.getMyAccount",
				UpdateMyAccount:     "server.requests.userAccount.updateMyAccount",
				BindGoogleAccount:   "server.requests.userAccount.bindGoogleAccount",
				UnbindGoogleAccount: "server.requests.userAccount.unbindGoogleAccount",
			},
			RootShelf: metricNameRootShelf{
				GetMyRootShelfById:        "server.requests.rootShelf.getMyRootShelfById",
				SearchRecentRootShelves:   "server.requests.rootShelf.searchRecentRootShelves",
				CreateRootShelf:           "server.requests.rootShelf.createRootShelf",
				CreateRootShelves:         "server.requests.rootShelf.createRootShelves",
				UpdateMyRootShelfById:     "server.requests.rootShelf.updateMyRootShelfById",
				UpdateMyRootShelvesByIds:  "server.requests.rootShelf.updateMyRootShelvesByIds",
				RestoreMyRootShelfById:    "server.requests.rootShelf.restoreMyRootShelfById",
				RestoreMyRootShelvesByIds: "server.requests.rootShelf.restoreMyRootShelvesByIds",
				DeleteMyRootShelfById:     "server.requests.rootShelf.deleteMyRootShelfById",
				DeleteMyRootShelvesByIds:  "server.requests.rootShelf.deleteMyRootShelvesByIds",
			},
			SubShelf: metricNameSubShelf{
				GetMySubShelfById:                       "server.requests.subShelf.getMySubShelfById",
				GetMySubShelvesByPrevSubShelfId:         "server.requests.subShelf.getMySubShelvesByPrevSubShelfId",
				GetAllMySubShelvesByRootShelfId:         "server.requests.subShelf.getAllMySubShelvesByRootShelfId",
				GetMySubShelvesAndItemsByPrevSubShelfId: "server.requests.subShelf.getMySubShelvesAndItemsByPrevSubShelfId",
				CreateSubShelfByRootShelfId:             "server.requests.subShelf.createSubShelfByRootShelfId",
				CreateSubShelvesByRootShelfIds:          "server.requests.subShelf.CreateSubShelvesByRootShelfIds",
				UpdateMySubShelfById:                    "server.requests.subShelf.updateMySubShelfById",
				UpdateMySubShelvesByIds:                 "server.requests.subShelf.UpdateMySubShelvesByIds",
				MoveMySubShelf:                          "server.requests.subShelf.moveMySubShelf",
				MoveMySubShelves:                        "server.requests.subShelf.moveMySubShelves",
				BatchMoveMySubShelves:                   "server.requests.subShelf.BatchMoveMySubShelves",
				RestoreMySubShelfById:                   "server.requests.subShelf.restoreMySubShelfById",
				RestoreMySubShelvesByIds:                "server.requests.subShelf.restoreMySubShelvesByIds",
				DeleteMySubShelfById:                    "server.requests.subShelf.deleteMySubShelfById",
				DeleteMySubShelvesByIds:                 "server.requests.subShelf.deleteMySubShelvesByIds",
			},
			Material: metricNameMaterial{
				GetMyMaterialById:                "server.requests.material.getMyMaterialById",
				GetMyMaterialAndItsParentById:    "server.requests.material.getMyMaterialAndItsParentById",
				GetMyMaterialsByParentSubShelfId: "server.requests.material.getMyMaterialsByParentSubShelfId",
				GetAllMyMaterialsByRootShelfId:   "server.requests.material.getAllMyMaterialsByRootShelfId",
				CreateTextbookMaterial:           "server.requests.material.createTextbookMaterial",
				CreateNotebookMaterial:           "server.requests.material.createNotebookMaterial",
				UpdateMyMaterialById:             "server.requests.material.updateMyMaterialById",
				SaveMyNotebookMaterialById:       "server.requests.material.saveMyNotebookMaterialById",
				MoveMyMaterialById:               "server.requests.material.moveMyMaterialById",
				MoveMyMaterialsByIds:             "server.requests.material.moveMyMaterialsByIds",
				RestoreMyMaterialById:            "server.requests.material.restoreMyMaterialById",
				RestoreMyMaterialsByIds:          "server.requests.material.restoreMyMaterialsByIds",
				DeleteMyMaterialById:             "server.requests.material.deleteMyMaterialById",
				DeleteMyMaterialsByIds:           "server.requests.material.deleteMyMaterialsByIds",
			},
			BlockPack: metricNameBlockPack{
				GetMyBlockPackById:                "server.requests.blockPack.getMyBlockPackById",
				GetMyBlockPackAndItsParentById:    "server.requests.blockPack.getMyBlockPackAndItsParentById",
				GetMyBlockPacksByParentSubShelfId: "server.requests.blockPack.getMyBlockPacksByParentSubShelfId",
				GetAllMyBlockPacksByRootShelfId:   "server.requests.blockPack.getAllMyBlockPacksByRootShelfId",
				CreateBlockPack:                   "server.requests.blockPack.createBlockPack",
				CreateBlockPacks:                  "server.requests.blockPack.CreateBlockPacks",
				UpdateMyBlockPackById:             "server.requests.blockPack.updateMyBlockPackById",
				UpdateMyBlockPacksByIds:           "server.requests.blockPack.UpdateMyBlockPacksByIds",
				MoveMyBlockPackById:               "server.requests.blockPack.moveMyBlockPackById",
				MoveMyBlockPacksByIds:             "server.requests.blockPack.moveMyBlockPacksByIds",
				BatchMoveMyBlockPacksByIds:        "server.requests.blockPack.BatchMoveMyBlockPacksByIds",
				RestoreMyBlockPackById:            "server.requests.blockPack.restoreMyBlockPackById",
				RestoreMyBlockPacksByIds:          "server.requests.blockPack.restoreMyBlockPacksByIds",
				DeleteMyBlockPackById:             "server.requests.blockPack.deleteMyBlockPackById",
				DeleteMyBlockPacksByIds:           "server.requests.blockPack.deleteMyBlockPacksByIds",
			},
			BlockGroup: metricNameBlockGroup{
				GetMyBlockGroupById:                                    "server.requests.blockGroup.getMyBlockGroupById",
				GetMyBlockGroupAndItsBlocksById:                        "server.requests.blockGroup.getMyBlockGroupAndItsBlocksById",
				GetMyBlockGroupsAndTheirBlocksByIds:                    "server.requests.blockGroup.getMyBlockGroupsAndTheirBlocksByIds",
				GetMyBlockGroupsAndTheirBlocksByBlockPackId:            "server.requests.blockGroup.getMyBlockGroupsAndTheirBlocksByBlockPackId",
				GetMyBlockGroupsByPrevBlockGroupId:                     "server.requests.blockGroup.getMyBlockGroupsByPrevBlockGroupId",
				GetAllMyBlockGroupsByBlockPackId:                       "server.requests.blockGroup.getAllMyBlockGroupsByBlockPackId",
				InsertBlockGroupByBlockPackId:                          "server.requests.blockGroup.insertBlockGroupByBlockPackId",
				InsertBlockGroupAndItsBlocksByBlockPackId:              "server.requests.blockGroup.insertBlockGroupAndItsBlocksByBlockPackId",
				InsertBlockGroupsAndTheirBlocksByBlockPackId:           "server.requests.blockGroup.insertBlockGroupsAndTheirBlocksByBlockPackId",
				InsertSequentialBlockGroupsAndTheirBlocksByBlockPackId: "server.requests.blockGroup.insertSequentialBlockGroupsAndTheirBlocksByBlockPackId",
				MoveMyBlockGroupsByIds:                                 "server.requests.blockGroup.moveMyBlockGroupsByIds",
				RestoreMyBlockGroupById:                                "server.requests.blockGroup.restoreMyBlockGroupById",
				RestoreMyBlockGroupsByIds:                              "server.requests.blockGroup.restoreMyBlockGroupsByIds",
				DeleteMyBlockGroupById:                                 "server.requests.blockGroup.deleteMyBlockGroupById",
				DeleteMyBlockGroupsByIds:                               "server.requests.blockGroup.deleteMyBlockGroupsByIds",
			},
			Block: metricNameBlock{
				GetMyBlockById:             "server.requests.block.getMyBlockById",
				GetMyBlocksByIds:           "server.requests.block.getMyBlocksByIds",
				GetMyBlocksByBlockGroupId:  "server.requests.block.getMyBlocksByBlockGroupId",
				GetMyBlocksByBlockGroupIds: "server.requests.block.getMyBlocksByBlockGroupIds",
				GetMyBlocksByBlockPackId:   "server.requests.block.getMyBlocksByBlockPackId",
				GetAllMyBlocks:             "server.requests.block.getAllMyBlocks",
				InsertBlock:                "server.requests.block.insertBlock",
				InsertBlocks:               "server.requests.block.insertBlocks",
				UpdateMyBlockById:          "server.requests.block.updateMyBlockById",
				UpdateMyBlocksByIds:        "server.requests.block.updateMyBlocksByIds",
				RestoreMyBlockById:         "server.requests.block.restoreMyBlockById",
				RestoreMyBlocksByIds:       "server.requests.block.restoreMyBlocksByIds",
				DeleteMyBlockById:          "server.requests.block.deleteMyBlockById",
				DeleteMyBlocksByIds:        "server.requests.block.deleteMyBlocksByIds",
			},
		},
		Responses: MetricNameResponse{
			Success: struct{ Total string }{
				Total: "server.responses.success.total",
			},
			Failed: struct {
				Total        string
				Timeout      string
				Unauthorized string
				RateLimit    string
			}{
				Total:        "server.responses.failed.total",
				Timeout:      "server.responses.failed.timeout",
				Unauthorized: "server.responses.failed.unauthorized",
				RateLimit:    "server.responses.failed.rateLimit",
			},
			Email: struct {
				Welcome  string
				AuthCode string
			}{
				Welcome:  "server.responses.email.welcome",
				AuthCode: "server.responses.email.authCode",
			},
		},
	},
}
