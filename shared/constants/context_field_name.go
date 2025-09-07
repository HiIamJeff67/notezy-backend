package constants

type ContextFieldName string

const (
	ContextFieldName_User_Id          ContextFieldName = "User-Id"
	ContextFieldName_User_PublicId    ContextFieldName = "User-PublicId"
	ContextFieldName_User_Name        ContextFieldName = "User-Name"
	ContextFieldName_User_DisplayName ContextFieldName = "User-DisplayName"
	ContextFieldName_User_Email       ContextFieldName = "User-Email"
	ContextFieldName_AccessToken      ContextFieldName = "AccessToken"
	ContextFieldName_User_Role        ContextFieldName = "User-Role"
	ContextFieldName_User_Plan        ContextFieldName = "User-Plan"

	ContextFieldName_GinContext          ContextFieldName = "GinContext"
	ContextFieldName_FormDataFileHeaders ContextFieldName = "FormDataFileHeaders"
)

func (cfn ContextFieldName) String() string {
	return string(cfn)
}
