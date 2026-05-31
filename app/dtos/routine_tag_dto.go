package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetOneRoutineTagByIdReqDto struct {
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

type CreateOneRoutineTagByUserIdReqDto struct {
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

type BulkCreateManyRoutineTagsByUserIdReqDto struct {
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

type UpdateOneRoutineTagByIdReqDto struct {
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

type BulkUpdateManyRoutineTagsByIdsReqDto struct {
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

type HardDeleteOneRoutineTagByIdReqDto struct {
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

type HardDeleteManyRoutineTagsByIdsReqDto struct {
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

type GetOneRoutineTagByIdResDto struct {
	Id        uuid.UUID            `json:"id"`
	Name      string               `json:"name"`
	Color     string               `json:"color"`
	Icon      *enums.SupportedIcon `json:"icon"`
	UpdatedAt time.Time            `json:"updatedAt"`
	CreatedAt time.Time            `json:"createdAt"`
}

type CreateOneRoutineTagByUserIdResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type BulkCreateManyRoutineTagsByUserIdResDto struct {
	Ids       []uuid.UUID `json:"ids"`
	CreatedAt time.Time   `json:"createdAt"`
}

type UpdateOneRoutineTagByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type BulkUpdateManyRoutineTagsByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type HardDeleteOneRoutineTagByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteManyRoutineTagsByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
