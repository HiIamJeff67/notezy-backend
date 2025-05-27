package models

import (
	"go-gorm-api/app/exceptions"
	"go-gorm-api/global"
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

/* ============================== Schema ============================== */
type User struct {
	Id				uuid.UUID		`json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	Name 			string			`json:"name" gorm:"column:name; not null; size:16;"`
	Email	 		string			`json:"email" gorm:"column:email; unique; not null;"`
	Password		string			`json:"password" gorm:"column:password; not null; size:255;"`
	RefreshToken 	string			`json:"refreshToken" gorm:"column:refresh_token; not null;"`
	Role			UserRole        `json:"role" gorm:"column:role; type:UserRole; not null; default:'Guest';"`
	Plan    		UserPlan		`json:"plan" gorm:"column:plan; type:UserPlan; not null; default:'Free';"`
	Status  		UserStatus		`json:"status" gorm:"column:status; type:UserStatus; not null; default:'Online';"`
	UpdatedAt       time.Time		`json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt 		time.Time		`json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relation
	UserInfo		UserInfo		`json:"userInfo" gorm:"foreignKey:UserId; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserAccount     UserAccount		`json:"userAccount" gorm:"foreignKey:UserId; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserSetting		UserSetting		`json:"userSetting" gorm:"foreignKey:UserId; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Badges 			[]Badge			`json:"badges" gorm:"many2many:\"UsersToBadgesTable\"; foreignKey:ID; joinForeignKey:UserID; references:ID; joinReferences:BadgeID;"`
}

// force gorm to use the given table name
func (User) TableName() string {
	return string(global.ValidTableName_UserTable)
}
/* ============================== Schema ============================== */

/* ============================== Input & Output ============================== */
type CreateUserInput struct {
	Name 			string			`json:"name" gorm:"column:name;"`
	Email	 		string			`json:"email" gorm:"column:email;"`
	Password		string			`json:"password" gorm:"column:password;"`
	RefreshToken	string			`json:"refreshToken" gorm:"column:refresh_token;"`
}
type UpdateUserInput struct {
	Name 			*string			`json:"name" gorm:"column:name;"`
	Email	 		*string			`json:"email" gorm:"column:email;"`
	Password		*string			`json:"password" gorm:"column:password;"`
	RefreshToken	*string			`json:"refreshToken" gorm:"column:refresh_token;"`
}
/* ============================== Input & Output ============================== */

/* ============================== Methods ============================== */
func GetUserById(id uuid.UUID) (user User, err error) {
	result := NotezyDB.Table(User{}.TableName()).Where("id = ?", id).First(&user)
	
	if err = result.Error; err != nil {
		return User{}, exceptions.User.NotFound().WithDetials(err).Log().Error
	}
	
	return user, nil
}

func GetAllUsers() (users []User, err error) {
	result := NotezyDB.Table(User{}.TableName()).Find(&users)
	return users, result.Error
}

func CreateUser(input CreateUserInput) (newUser User, err error) {
	newUser = User{
		Name: input.Name, 
		Email: input.Email,
		Password: input.Password,
		RefreshToken: input.RefreshToken,
	}
	result := NotezyDB.Table(User{}.TableName()).Create(&newUser)
	if err = result.Error; err != nil {
		return User{}, exceptions.User.FailedToCreate().WithDetials(err).Log().Error
	}
	return newUser, nil
}

func UpdateUserById(id uuid.UUID, input UpdateUserInput) (updatedUser User, err error) {
	updatedUser = User{
		Name: *input.Name, 
		Email: *input.Email,
		Password: *input.Password,
		RefreshToken: *input.RefreshToken,
	}

	result := NotezyDB.Table(User{}.TableName()).Where("id = ?", id).Updates(&updatedUser)

	if err = result.Error; err != nil {
		return User{}, exceptions.User.FailedToUpdate().WithDetials(err).Log().Error
	}
	return updatedUser, nil
}

func DeleteUserById(id uuid.UUID) (deletedUser User, err error) {
	tx := NotezyDB.Begin()

	result := tx.Table(User{}.TableName()).Where("id = ?", id).First(&deletedUser)
	if err = result.Error; err != nil {
		tx.Rollback()
		return User{}, err
	}

	result = tx.Table(User{}.TableName()).Delete(&deletedUser)
	if err = result.Error; err != nil {
		tx.Rollback()
		return User{}, exceptions.User.FailedToDelete().WithDetials(err).Log().Error
	}

	return deletedUser, tx.Commit().Error
}
/* ============================== Methods ============================== */