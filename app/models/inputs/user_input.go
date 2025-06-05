package inputs

import "notezy-backend/app/models/enums"

type CreateUserInput struct {
	Name         string  `json:"name" validate:"required,min=6,max=16,alphanum" gorm:"column:name;"`
	DisplayName  string  `json:"displayName" validate:"required,min=6,max=32" gorm:"column:display_name"`
	Email        string  `json:"email" validate:"required,email" gorm:"column:email;"`
	Password     string  `json:"password" validate:"required,min=8,max=1024" gorm:"column:password;"` // hashed password
	RefreshToken *string `json:"refreshToken" validate:"omitempty" gorm:"column:refresh_token;"`
	UserAgent    string  `json:"userAgent" validate:"required" gorm:"column:user_agent;"`
}
type UpdateUserInput struct {
	Name         *string           `json:"name" validate:"omitempty,min=6,max=16,alphanum" gorm:"column:name;"`
	DisplayName  *string           `json:"displayName" validate:"omitempty,min=6,max=32,alphanum" gorm:"column:display_name;"`
	Email        *string           `json:"email" validate:"omitempty,email" gorm:"column:email;"`
	Password     *string           `json:"password" validate:"omitempty,min=8,max=1024" gorm:"column:password;"` // hashed password
	RefreshToken *string           `json:"refreshToken" validate:"omitempty" gorm:"column:refresh_token;"`
	UserAgent    *string           `json:"userAgent" validate:"omitempty" gorm:"column:user_agent;"`
	Role         *enums.UserRole   `json:"role" validate:"omitempty,isrole" gorm:"column:role;"`
	Plan         *enums.UserPlan   `json:"plan" validate:"omitempty,isplan" gorm:"column:plan;"`
	Status       *enums.UserStatus `json:"status" validate:"omitempty,isstatus" gorm:"column:status;"`
}
