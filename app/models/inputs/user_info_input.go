package inputs

import (
	"time"

	enums "notezy-backend/app/models/schemas/enums"
)

type CreateUserInfoInput struct {
	CoverBackgroundURL *string           `json:"coverBackgroundURL" validate:"omitempty" gorm:"column:cover_background_url;"`
	AvatarURL          *string           `json:"avatarURL" validate:"omitempty" gorm:"column:avatar_url;"`
	Header             *string           `json:"header" validate:"omitempty,min=0,max=64" gorm:"column:header;"`
	Introduction       *string           `json:"introduction" validate:"omitempty,min=0,max=256" gorm:"column:introduction;"`
	Gender             *enums.UserGender `json:"gender" validate:"omitempty,isgender" gorm:"column:gender;"`
	Country            *enums.Country    `json:"country" validate:"omitempty,iscountry" gorm:"column:country;"`
	BirthDate          *time.Time        `json:"birthDate" validate:"omitempty" gorm:"column:birth_date;"`
}

type UpdateUserInfoInput struct {
	CoverBackgroundURL *string           `json:"coverBackgroundURL" validate:"omitempty" gorm:"column:cover_background_url;"`
	AvatarURL          *string           `json:"avatarURL" validate:"omitempty" gorm:"column:avatar_url;"`
	Header             *string           `json:"header" validate:"omitempty,min=0,max=64" gorm:"column:header;"`
	Introduction       *string           `json:"introduction" validate:"omitempty,min=0,max=256" gorm:"column:introduction;"`
	Gender             *enums.UserGender `json:"gender" validate:"omitempty,isgender" gorm:"column:gender;"`
	Country            *enums.Country    `json:"country" validate:"omitempty,iscountry" gorm:"column:country;"`
	BirthDate          *time.Time        `json:"birthDate" validate:"omitempty" gorm:"column:birth_date;"`
}

type PartialUpdateUserInfoInput = PartialUpdateInput[UpdateUserInfoInput]
