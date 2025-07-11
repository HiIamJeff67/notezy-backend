package inputs

import (
	enums "notezy-backend/app/models/schemas/enums"
)

type CreateUserSettingInput struct {
	Language           *enums.Language `json:"language" validate:"omitempty,islanguage" gorm:"column:language;"`
	GeneralSettingCode *int            `json:"generalSettingCode" validate:"omitempty,min=0,max=999999999" gorm:"column:general_setting_code;"`
	PrivacySettingCode *int            `json:"privacySettingCode" validate:"omitempty,min=0,max=999999999" gorm:"column:privacy_setting_code;"`
}

type UpdateUserSettingInput struct {
	Language           *enums.Language `json:"language" validate:"omitempty,islanguage" gorm:"column:language;"`
	GeneralSettingCode *int            `json:"generalSettingCode" validate:"omitempty,min=0,max=999999999" gorm:"column:general_setting_code;"`
	PrivacySettingCode *int            `json:"privacySettingCode" validate:"omitempty,min=0,max=999999999" gorm:"column:privacy_setting_code;"`
}

type PartialUpdateUserSettingInput = PartialUpdateInput[UpdateUserSettingInput]
