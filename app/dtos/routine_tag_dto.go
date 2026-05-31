package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetMyRoutineTagByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			RoutineTagId uuid.UUID `form:"routineTagId" validate:"required"`
		},
	]
}

type CreateRoutineTagReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			Id    *uuid.UUID           `json:"id" validate:"omitnil"`
			Name  string               `json:"name" validate:"required,min=1,max=128"`
			Color string               `json:"color" validate:"omitempty,ishexcodecolor"`
			Icon  *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
		},
		any,
	]
}

type CreateRoutineTagsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			CreatedRoutineTags []struct {
				Id    *uuid.UUID           `json:"id" validate:"omitnil"`
				Name  string               `json:"name" validate:"required,min=1,max=128"`
				Color string               `json:"color" validate:"omitempty,ishexcodecolor"`
				Icon  *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
			} `json:"createdRoutineTags" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type UpdateMyRoutineTagByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
			PartialUpdateDto[struct {
				Name  *string              `json:"name" validate:"omitnil,min=1,max=128"`
				Color *string              `json:"color" validate:"omitnil,ishexcodecolor"`
				Icon  *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
			}]
		},
		any,
	]
}

type UpdateMyRoutineTagsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			UpdatedRoutineTags []struct {
				RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
				PartialUpdateDto[struct {
					Name  *string              `json:"name" validate:"omitnil,min=1,max=128"`
					Color *string              `json:"color" validate:"omitnil,ishexcodecolor"`
					Icon  *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
				}]
			} `json:"updatedRoutineTags" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type HardDeleteMyRoutineTagByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
		},
		any,
	]
}

type HardDeleteMyRoutineTagsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineTagIds []uuid.UUID `json:"routineTagIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyRoutineTagByIdResDto struct {
	Id        uuid.UUID            `json:"id"`
	Name      string               `json:"name"`
	Color     string               `json:"color"`
	Icon      *enums.SupportedIcon `json:"icon"`
	UpdatedAt time.Time            `json:"updatedAt"`
	CreatedAt time.Time            `json:"createdAt"`
}

type CreateRoutineTagResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateRoutineTagsResDto struct {
	Ids       []uuid.UUID `json:"ids"`
	CreatedAt time.Time   `json:"createdAt"`
}

type UpdateMyRoutineTagByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMyRoutineTagsByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type HardDeleteMyRoutineTagByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteMyRoutineTagsByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
