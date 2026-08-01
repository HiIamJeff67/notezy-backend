package durablejobexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

type userSettingExceptionDomain struct{}

var UserSetting = userSettingExceptionDomain{}

func (userSettingExceptionDomain) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"UserSetting",
		"Repository",
		"UserSetting was not found",
		http.StatusNotFound,
	)
}

func (userSettingExceptionDomain) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"UserSetting",
		"Repository",
		"Failed to create UserSetting",
		http.StatusInternalServerError,
		true,
	)
}

func (userSettingExceptionDomain) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"UserSetting",
		"Repository",
		"Failed to update UserSetting",
		http.StatusInternalServerError,
		true,
	)
}

func (userSettingExceptionDomain) FailedToDelete() *exceptions.Exception {
	return exceptions.New(
		"FailedToDelete",
		"UserSetting",
		"Repository",
		"Failed to delete UserSetting",
		http.StatusInternalServerError,
		true,
	)
}

func (userSettingExceptionDomain) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"UserSetting",
		"Repository",
		"No changes were applied to UserSetting",
		http.StatusNotModified,
	)
}
