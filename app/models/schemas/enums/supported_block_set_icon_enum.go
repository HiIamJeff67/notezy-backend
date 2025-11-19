package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

/* ============================== Definition ============================== */

type SupportedBlockPackIcon string

const (
	SupportedBlockPackIcon_GrinningFace               SupportedBlockPackIcon = "😀"
	SupportedBlockPackIcon_SmilingFaceWithSmilingEyes SupportedBlockPackIcon = "😊"
	SupportedBlockPackIcon_RedHeart                   SupportedBlockPackIcon = "❤️"
	SupportedBlockPackIcon_Fire                       SupportedBlockPackIcon = "🔥"
	SupportedBlockPackIcon_Star                       SupportedBlockPackIcon = "⭐"
	SupportedBlockPackIcon_Books                      SupportedBlockPackIcon = "📚"
	SupportedBlockPackIcon_Notebook                   SupportedBlockPackIcon = "📓"
	SupportedBlockPackIcon_PencilPaper                SupportedBlockPackIcon = "📝"
	SupportedBlockPackIcon_Lightbulb                  SupportedBlockPackIcon = "💡"
	SupportedBlockPackIcon_Rocket                     SupportedBlockPackIcon = "🚀"
	SupportedBlockPackIcon_CheckMark                  SupportedBlockPackIcon = "✅"
	SupportedBlockPackIcon_Pin                        SupportedBlockPackIcon = "📌"
	SupportedBlockPackIcon_FolderOpen                 SupportedBlockPackIcon = "📂"
	SupportedBlockPackIcon_Calendar                   SupportedBlockPackIcon = "📅"
	SupportedBlockPackIcon_Clock                      SupportedBlockPackIcon = "⏰"
)

/* ============================== All Instances ============================== */

var AllSupportedBlockPackIcons = []SupportedBlockPackIcon{
	SupportedBlockPackIcon_GrinningFace,
	SupportedBlockPackIcon_SmilingFaceWithSmilingEyes,
	SupportedBlockPackIcon_RedHeart,
	SupportedBlockPackIcon_Fire,
	SupportedBlockPackIcon_Star,
	SupportedBlockPackIcon_Books,
	SupportedBlockPackIcon_Notebook,
	SupportedBlockPackIcon_PencilPaper,
	SupportedBlockPackIcon_Lightbulb,
	SupportedBlockPackIcon_Rocket,
	SupportedBlockPackIcon_CheckMark,
	SupportedBlockPackIcon_Pin,
	SupportedBlockPackIcon_FolderOpen,
	SupportedBlockPackIcon_Calendar,
	SupportedBlockPackIcon_Clock,
}

var AllSupportedBlockPackIconStrings = []string{
	string(SupportedBlockPackIcon_GrinningFace),
	string(SupportedBlockPackIcon_SmilingFaceWithSmilingEyes),
	string(SupportedBlockPackIcon_RedHeart),
	string(SupportedBlockPackIcon_Fire),
	string(SupportedBlockPackIcon_Star),
	string(SupportedBlockPackIcon_Books),
	string(SupportedBlockPackIcon_Notebook),
	string(SupportedBlockPackIcon_PencilPaper),
	string(SupportedBlockPackIcon_Lightbulb),
	string(SupportedBlockPackIcon_Rocket),
	string(SupportedBlockPackIcon_CheckMark),
	string(SupportedBlockPackIcon_Pin),
	string(SupportedBlockPackIcon_FolderOpen),
	string(SupportedBlockPackIcon_Calendar),
	string(SupportedBlockPackIcon_Clock),
}

/* ============================== Methods ============================== */

func (bssi SupportedBlockPackIcon) Name() string {
	return reflect.TypeOf(bssi).Name()
}

func (bssi *SupportedBlockPackIcon) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*bssi = SupportedBlockPackIcon(string(v))
		return nil
	case string:
		*bssi = SupportedBlockPackIcon(v)
		return nil
	}
	return scanError(value, bssi)
}

func (bssi SupportedBlockPackIcon) Value() (driver.Value, error) {
	return string(bssi), nil
}

func (bssi SupportedBlockPackIcon) String() string {
	return string(bssi)
}

func (bssi *SupportedBlockPackIcon) IsValidEnum() bool {
	return slices.Contains(AllSupportedBlockPackIcons, *bssi)
}

func ConvertStringToSupportedBlockPackIcon(enumString string) (*SupportedBlockPackIcon, error) {
	for _, supportedBlockPackIcon := range AllSupportedBlockPackIcons {
		if string(supportedBlockPackIcon) == enumString {
			return &supportedBlockPackIcon, nil
		}
	}
	return nil, fmt.Errorf("invalid access control permission: %s", enumString)
}
