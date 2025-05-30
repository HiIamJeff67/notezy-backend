package models

import (
	"time"

	"gorm.io/gorm"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"

	exceptions "notezy-backend/app/exceptions"
	global "notezy-backend/global"
)

/* ============================== Schema ============================== */
type User struct {
	Id           uuid.UUID  `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	Name         string     `json:"name" gorm:"column:name; unique; not null; size:16;"`
	DisplayName  string     `json:"displayName" gorm:"column:display_name; not null; size:32;"`
	Email        string     `json:"email" gorm:"column:email; unique; not null;"`
	Password     string     `json:"password" gorm:"column:password; not null; size:255;"`
	RefreshToken string     `json:"refreshToken" gorm:"column:refresh_token; not null;"`
	Role         UserRole   `json:"role" gorm:"column:role; type:UserRole; not null; default:'Guest';"`
	Plan         UserPlan   `json:"plan" gorm:"column:plan; type:UserPlan; not null; default:'Free';"`
	Status       UserStatus `json:"status" gorm:"column:status; type:UserStatus; not null; default:'Online';"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relation
	UserInfo    UserInfo    `json:"userInfo" gorm:"foreignKey:UserId; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserAccount UserAccount `json:"userAccount" gorm:"foreignKey:UserId; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserSetting UserSetting `json:"userSetting" gorm:"foreignKey:UserId; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Badges      []Badge     `json:"badges" gorm:"many2many:\"UsersToBadgesTable\"; foreignKey:ID; joinForeignKey:UserID; references:ID; joinReferences:BadgeID;"`
}

// force gorm to use the given table name
func (User) TableName() string {
	return string(global.ValidTableName_UserTable)
}

/* ============================== Input & Output ============================== */
type CreateUserInput struct {
	Name         string `json:"name" validate:"required, min=6, max=16, alphanum" gorm:"column:name;"`
	Email        string `json:"email" validate:"required, email" gorm:"column:email;"`
	Password     string `json:"password" validate:"required, min=8, max=32, isstrongpassword" gorm:"column:password;"`
	RefreshToken string `json:"refreshToken" gorm:"column:refresh_token;"`
}
type UpdateUserInput struct {
	Name         *string    `json:"name" validate:"required, min=6, max=16, alphanum" gorm:"column:name;"`
	DisplayName  string     `json:"displayName" validae:"min=6, max=32, alphanum" gorm:"column:display_name; not null; size:32;"`
	Email        *string    `json:"email" validate:"required, email" gorm:"column:email;"`
	Password     *string    `json:"password" validate:"required, min=8, max=32, isstrongpassword" gorm:"column:password;"`
	RefreshToken *string    `json:"refreshToken" gorm:"column:refresh_token;"`
	Role         UserRole   `json:"role" validate:"oneof=Admin Normal Guest" gorm:"column:role; type:UserRole; not null; default:'Guest';"`
	Plan         UserPlan   `json:"plan" gorm:"column:plan; type:UserPlan; not null; default:'Free';"`
	Status       UserStatus `json:"status" gorm:"column:status; type:UserStatus; not null; default:'Online';"`
}

/* ============================== Methods ============================== */
func GetUserById(db *gorm.DB, id uuid.UUID) (User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	user := User{}
	result := db.Table(User{}.TableName()).Where("id = ?", id).First(&user)

	if err := result.Error; err != nil {
		return User{}, exceptions.User.NotFound().WithError(err)
	}

	return user, nil
}

func GetAllUsers(db *gorm.DB) ([]User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}
	users := []User{}
	result := db.Table(User{}.TableName()).Find(&users)
	return users, exceptions.User.NotFound().WithError(result.Error)
}

func CreateUser(db *gorm.DB, input CreateUserInput) (User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	if err := Validator.Struct(input); err != nil {
		return User{}, exceptions.User.InvalidInput().WithError(err)
	}

	newUser := User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     input.Password,
		RefreshToken: input.RefreshToken,
	}
	result := db.Table(User{}.TableName()).Create(&newUser)
	if err := result.Error; err != nil {
		return User{}, exceptions.User.FailedToCreate().WithError(err)
	}
	return newUser, nil
}

func UpdateUserById(db *gorm.DB, id uuid.UUID, input UpdateUserInput) (User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	if err := Validator.Struct(input); err != nil {
		return User{}, exceptions.User.InvalidInput().WithError(err)
	}

	updatedUser := User{
		Name:         *input.Name,
		Email:        *input.Email,
		Password:     *input.Password,
		RefreshToken: *input.RefreshToken,
	}

	result := db.Table(User{}.TableName()).Where("id = ?", id).Updates(&updatedUser)

	if err := result.Error; err != nil {
		return User{}, exceptions.User.FailedToUpdate().WithError(err)
	}
	return updatedUser, nil
}

func DeleteUserById(db *gorm.DB, id uuid.UUID) (User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	tx := db.Begin()

	deletedUser := User{}
	result := tx.Table(User{}.TableName()).Where("id = ?", id).First(&deletedUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return User{}, exceptions.User.NotFound().WithError(err)
	}

	result = tx.Table(User{}.TableName()).Delete(&deletedUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return User{}, exceptions.User.FailedToDelete().WithError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return User{}, exceptions.User.FailedToDelete().WithError(err)
	}

	return deletedUser, nil
}
