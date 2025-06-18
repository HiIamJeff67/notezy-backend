package inputs

type CreateThemeInput struct {
	Name        string `json:"name" validate:"required,min=1,max=100" gorm:"column:name;"`
	Version     string `json:"version" validate:"required,max=20" gorm:"column:is_default;"`
	DownloadURL string `json:"downloadURL" validate:"required,url" gorm:"column:download_url;"`
	IsDefault   bool   `json:"isDefault" validate:"required" gorm:"column:is_default;"`
}

type UpdateThemeInput struct {
	Name        *string `json:"name" validate:"required,min=1,max=100" gorm:"column:name;"`
	Version     *string `json:"version" validate:"required,max=20" gorm:"column:is_default;"`
	DownloadURL *string `json:"downloadURL" validate:"required,url" gorm:"column:download_url;"`
	IsDefault   *bool   `json:"isDefault" validate:"required" gorm:"column:is_default;"`
}

type PartialUpdateThemeInput = PartialUpdateInput[UpdateThemeInput]
