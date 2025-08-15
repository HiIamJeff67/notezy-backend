package constants

type ContextFieldName string

const (
	ContextFieldName_User_Id          ContextFieldName = "user-Id"
	ContextFieldName_User_PublicId    ContextFieldName = "user-PublicId"
	ContextFieldName_User_Name        ContextFieldName = "user-Name"
	ContextFieldName_User_DisplayName ContextFieldName = "user-DisplayName"
	ContextFieldName_User_Email       ContextFieldName = "user-Email"
	ContextFieldName_AccessToken      ContextFieldName = "accessToken"
	ContextFieldName_User_Role        ContextFieldName = "user-Role"
	ContextFieldName_User_Plan        ContextFieldName = "user-plan"
	ContextFieldName_Gin_Context      ContextFieldName = "gin-Context"
)

func (cfn ContextFieldName) String() string {
	return string(cfn)
}
