package inputs

import (
	enums "notezy-backend/app/models/enums"
)

type CreateUserSettingInput struct {
	Language           *enums.Language `json:"language" validate:"islanguage" gorm:"column:language;"`
	GeneralSettingCode *int            `json:"generalSettingCode" validate:"min=0,max=999999999" gorm:"column:general_setting_code;"`
	PrivacySettingCode *int            `json:"privacySettingCode" validate:"min=0,max=999999999" gorm:"column:privacy_setting_code;"`
}

type UpdateUserSettingInput struct {
	Language           *enums.Language `json:"language" validate:"islanguage" gorm:"column:language;"`
	GeneralSettingCode *int            `json:"generalSettingCode" validate:"min=0,max=999999999" gorm:"column:general_setting_code;"`
	PrivacySettingCode *int            `json:"privacySettingCode" validate:"min=0,max=999999999" gorm:"column:privacy_setting_code;"`
}

type PartialUpdateUserSettingInput = PartialUpdateInput[UpdateUserSettingInput]
