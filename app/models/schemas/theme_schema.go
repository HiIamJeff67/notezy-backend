package schemas

import (
	"time"

	"github.com/google/uuid"

	shared "notezy-backend/app/shared"
)

type Theme struct {
	Id            uuid.UUID `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	Name          string    `json:"name" gorm:"column:name; unique; not null;"`
	AuthorId      uuid.UUID `json:"authorId" gorm:"column:author_id; type:uuid; not null; uniqueIndex;"`
	Version       string    `json:"version" gorm:"column:version; not null;"`
	IsDefault     bool      `json:"isDefault" gorm:"column:is_default; type:boolean; not null;"`
	DownloadURL   *string   `json:"downloadURL" gorm:"column:download_url;"`
	DownloadCount int64     `json:"downloadCount" gorm:"column:download_count; type:bigint; not null default:'0';"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`

	// relation
	Author User `json:"author" gorm:"foreignKey:AuthorId; references:Id;"`
}

func (Theme) TableName() string {
	return shared.ValidTableName_ThemeTable.String()
}
