package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
)

type UserSettingException struct {
	CoreException
}

func NewUserSettingException() UserSettingException {
	return UserSettingException{CoreException: NewCoreException("UserSetting")}
}

func (UserSettingException) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"UserSetting",
		"Repository",
		"UserSetting was not found",
		http.StatusNotFound,
	)
}

func (UserSettingException) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"UserSetting",
		"Repository",
		"Failed to create UserSetting",
		http.StatusInternalServerError,
		true,
	)
}

func (UserSettingException) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"UserSetting",
		"Repository",
		"Failed to update UserSetting",
		http.StatusInternalServerError,
		true,
	)
}

func (UserSettingException) FailedToDelete() *exceptions.Exception {
	return exceptions.New(
		"FailedToDelete",
		"UserSetting",
		"Repository",
		"Failed to delete UserSetting",
		http.StatusInternalServerError,
		true,
	)
}

func (UserSettingException) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"UserSetting",
		"Repository",
		"No changes were applied to UserSetting",
		http.StatusNotModified,
	)
}
