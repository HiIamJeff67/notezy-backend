package models

import (
	"go-gorm-api/global"
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

/* ============================== Schema ============================== */
type Badge struct {
	Id				uuid.UUID		`json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid()"`
	Title			string			`json:"title" gorm:"column:title; not null; size:64;"`
	Description     string			`json:"description" gorm:"column:description; not null; size:256;"`
	Type  			BadgeType		`json:"type" gorm:"column:type; type:BadgeType; not null; default:'Brozen';"`
	ImageURL		string			`json:"imageURL" gorm:"column:image_url;"`
	CreatedAt		time.Time		`json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relation
	Users			[]User 			`json:"users" gorm:"many2many:\"UsersToBadgesTable\"; foreignKey:ID; joinForeignKey:BadgeID; references:ID; joinReferences:UserID;"`
}

func (Badge) TableName() string {
	return string(global.ValidTableName_BadgeTable)
}
/* ============================== Schema ============================== */

/* ============================== Input ============================== */

/* ============================== Input ============================== */

/* ============================== Methods ============================== */

/* ============================== Methods ============================== */