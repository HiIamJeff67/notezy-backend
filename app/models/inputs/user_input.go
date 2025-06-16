package inputs

import (
	"time"

	enums "notezy-backend/app/models/enums"
)

type CreateUserInput struct {
	Name         string `json:"name" validate:"required,min=6,max=16,alphanum" gorm:"column:name;"`
	DisplayName  string `json:"displayName" validate:"required,min=6,max=32" gorm:"column:display_name"`
	Email        string `json:"email" validate:"required,email" gorm:"column:email;"`
	Password     string `json:"password" validate:"required,min=8,max=1024" gorm:"column:password;"` // hashed password
	RefreshToken string `json:"refreshToken" validate:"omitempty" gorm:"column:refresh_token;"`
	UserAgent    string `json:"userAgent" validate:"required" gorm:"column:user_agent;"`
}

type UpdateUserInput struct {
	Name           *string           `json:"name" gorm:"column:name;"`                        // validate:"omitempty,min=6,max=16,alphanum"
	DisplayName    *string           `json:"displayName" gorm:"column:display_name;"`         // validate:"omitempty,min=6,max=32,alphanum"
	Email          *string           `json:"email" gorm:"column:email;"`                      //validate:"omitempty,email"
	Password       *string           `json:"password" gorm:"column:password;"`                // validate:"omitempty,min=8,max=1024"
	RefreshToken   *string           `json:"refreshToken" gorm:"column:refresh_token;"`       // validate:"omitempty"
	LoginCount     *int32            `json:"loginCount" gorm:"column:login_count;"`           // validate:"omitempty,min=0"
	BlockLoginUtil *time.Time        `json:"BlockLoginUntil" gorm:"column:block_login_until"` // validate:"omitempty"
	UserAgent      *string           `json:"userAgent" gorm:"column:user_agent;"`             // validate:"omitempty"
	Role           *enums.UserRole   `json:"role" gorm:"column:role;"`                        // validate:"omitempty,isrole"
	Plan           *enums.UserPlan   `json:"plan" gorm:"column:plan;"`                        // validate:"omitempty,isplan"
	Status         *enums.UserStatus `json:"status" gorm:"column:status;"`                    // validate:"omitempty,isstatus"
}

type PartialUpdateUserInput = PartialUpdateInput[UpdateUserInput]

type DeleteUserInput struct {
	Name     string `json:"name" validate:"required,min=6,max=16,alphanum" gorm:"column:name;"`
	Password string `json:"password" validate:"required,min=8,max=1024" gorm:"column:password;"`
}
