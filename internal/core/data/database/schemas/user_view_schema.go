package schemas

// UserView is intentionally not an active GORM model yet.
//
// Core owns UserTable and the corresponding PostgreSQL view is created by
// schemas/views/user_view.sql. Notification currently uses a different
// PostgreSQL database, so it cannot query this view directly. Keeping this
// design here documents the future read-only projection without making a
// misleading cross-database model available to another runtime.
//
// When a runtime shares Core's database connection, uncomment the model below
// and use it only for reads. Its fields must remain limited to the data that
// the consuming runtime actually needs; permissions are a separate concern.
//
// type UserView struct {
//	PublicId uuid.UUID        `json:"publicId" gorm:"column:public_id;type:uuid;primaryKey;->"`
//	Status   enums.UserStatus `json:"status" gorm:"column:status;type:\"UserStatus\";->"`
// }
//
// func (UserView) TableName() string {
//	return "UserView"
// }
