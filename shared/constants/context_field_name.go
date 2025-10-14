package constants

type ContextFieldName string

// use "-" between value to represent the relationship or domain
// ex. "User-Id" means the id of the user, since we may use "Other-Id" to represent the id of some other stuff

const (
	ContextFieldName_User_Id          ContextFieldName = "User-Id"          // UUID
	ContextFieldName_User_PublicId    ContextFieldName = "User-PublicId"    // UUID
	ContextFieldName_User_Name        ContextFieldName = "User-Name"        // string
	ContextFieldName_User_DisplayName ContextFieldName = "User-DisplayName" // string
	ContextFieldName_User_Email       ContextFieldName = "User-Email"       // string
	ContextFieldName_IsNewAccessToken ContextFieldName = "IsNewAccessToken" // bool
	ContextFieldName_AccessToken      ContextFieldName = "AccessToken"      // string
	ContextFieldName_User_Role        ContextFieldName = "User-Role"        // enums.UserRole
	ContextFieldName_User_Plan        ContextFieldName = "User-Plan"        // enums.UserPlan

	ContextFieldName_GinContext          ContextFieldName = "GinContext"          // gin.Context
	ContextFieldName_FormDataFileHeaders ContextFieldName = "FormDataFileHeaders" // []*multipart.FileHeader
)

func (cfn ContextFieldName) String() string {
	return string(cfn)
}
