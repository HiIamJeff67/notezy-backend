package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	uuid "github.com/google/uuid"

	exceptions "notezy-backend/app/exceptions"
	"notezy-backend/app/util"
	global "notezy-backend/global"
)

/* ============================== Schema ============================== */
type User struct {
	Id           uuid.UUID  `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	Name         string     `json:"name" gorm:"column:name; unique; not null; size:16;"`
	DisplayName  string     `json:"displayName" gorm:"column:display_name; not null; size:32;"`
	Email        string     `json:"email" gorm:"column:email; unique; not null;"`
	Password     string     `json:"password" gorm:"column:password; not null; size:1024;"` // since we store the hashed password which is quite long
	RefreshToken string     `json:"refreshToken" gorm:"column:refresh_token; not null; default:'';"`
	Role         UserRole   `json:"role" gorm:"column:role; type:UserRole; not null; default:'Guest';"`
	Plan         UserPlan   `json:"plan" gorm:"column:plan; type:UserPlan; not null; default:'Free';"`
	Status       UserStatus `json:"status" gorm:"column:status; type:UserStatus; not null; default:'Online';"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relation
	UserInfo    UserInfo    `json:"userInfo" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserAccount UserAccount `json:"userAccount" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserSetting UserSetting `json:"userSetting" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Badges      []Badge     `json:"badges" gorm:"-"` // many2many:\"UsersToBadgesTable\"; foreignKey:Id; joinForeignKey:UserId; references:Id; joinReferences:BadgeId;
}

// force gorm to use the given table name
func (User) TableName() string {
	return string(global.ValidTableName_UserTable)
}

/* ============================== Input & Output ============================== */
type CreateUserInput struct {
	Name         string  `json:"name" validate:"required,min=6,max=16,alphanum" gorm:"column:name;"`
	DisplayName  string  `json:"displayName" validate:"required,min=6,max=32" gorm:"column:display_name"`
	Email        string  `json:"email" validate:"required,email" gorm:"column:email;"`
	Password     string  `json:"password" validate:"required,min=8,max=1024" gorm:"column:password;"`
	RefreshToken *string `json:"refreshToken" validate:"omitempty" gorm:"column:refresh_token;"`
}
type UpdateUserInput struct {
	Name         *string     `json:"name" validate:"omitempty,min=6,max=16,alphanum" gorm:"column:name;"`
	DisplayName  *string     `json:"displayName" validae:"omitempty,min=6,max=32,alphanum" gorm:"column:display_name;"`
	Email        *string     `json:"email" validate:"omitempty,email" gorm:"column:email;"`
	Password     *string     `json:"password" validate:"omitempty,min=8,max=1024" gorm:"column:password;"`
	RefreshToken *string     `json:"refreshToken" validate:"omitempty" gorm:"column:refresh_token;"`
	Role         *UserRole   `json:"role" validate:"omitempty,isrole" gorm:"column:role;"`
	Plan         *UserPlan   `json:"plan" validate:"omitempty,isplan" gorm:"column:plan;"`
	Status       *UserStatus `json:"status" validate:"omitempty,isstatus" gorm:"column:status;"`
}

/* ============================== Methods ============================== */
func GetUserById(db *gorm.DB, id uuid.UUID) (*User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	user := User{}
	result := db.Table(User{}.TableName()).
		Where("id = ?", id).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func GetUserByName(db *gorm.DB, name string) (*User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	user := User{}
	result := db.Table(User{}.TableName()).
		Where("name = ?", name).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	user := User{}
	result := db.Table(User{}.TableName()).
		Where("email = ?", email).
		First(&user)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	return &user, nil
}

func GetAllUsers(db *gorm.DB) (*[]User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	users := []User{}
	result := db.Table(User{}.TableName()).
		Find(&users)
	if err := result.Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(result.Error)
	}
	return &users, nil
}

func CreateUser(db *gorm.DB, input CreateUserInput) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	if err := Validator.Struct(input); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	// note that the create operation in gorm will NOT return anything
	// but the default value we set in gorm field in the above struct will be returned if we specified it in the "returning"
	var newUser User
	util.CopyNonNilFields(&newUser, input)
	result := db.Table(User{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"}, // for the following procedure such as create user info, create user account, generate refresh token etc..
		}}).
		Create(&newUser)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToCreate().WithError(err)
	}
	return &newUser.Id, nil
}

func UpdateUserById(db *gorm.DB, id uuid.UUID, input UpdateUserInput) (*User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	if err := Validator.Struct(input); err != nil {
		return nil, exceptions.User.InvalidInput().WithError(err)
	}

	var updatedUser User
	updatedUser.UpdatedAt = time.Now()
	util.CopyNonNilFields(&updatedUser, input)
	result := db.Table(User{}.TableName()).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Updates(&updatedUser)
	if err := result.Error; err != nil {
		return nil, exceptions.User.FailedToUpdate().WithError(err)
	}
	return &updatedUser, nil
}

func DeleteUserById(db *gorm.DB, id uuid.UUID) (*User, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	tx := db.Begin()

	deletedUser := User{}
	result := tx.Table(User{}.TableName()).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		First(&deletedUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, exceptions.User.NotFound().WithError(err)
	}

	result = tx.Table(User{}.TableName()).
		Delete(&deletedUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, exceptions.User.FailedToDelete().WithError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.User.FailedToDelete().WithError(err)
	}

	return &deletedUser, nil
}
