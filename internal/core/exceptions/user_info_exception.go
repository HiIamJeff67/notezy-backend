package apiexceptions

type UserInfoException struct {
	CoreException
}

func NewUserInfoException() UserInfoException {
	return UserInfoException{
		CoreException: NewCoreException("UserInfo"),
	}
}
