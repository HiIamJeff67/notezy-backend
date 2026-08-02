package blocksdto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type GetMyBlockByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			BlockId uuid.UUID `json:"blockId" validate:"required"`
		},
		struct{},
	]
}

type BlockResponseDto struct {
	Id            uuid.UUID      `json:"id"`
	BlockPackId   uuid.UUID      `json:"blockPackId"`
	ParentBlockId *uuid.UUID     `json:"parentBlockId"`
	PrevBlockId   *uuid.UUID     `json:"prevBlockId"`
	NextBlockId   *uuid.UUID     `json:"nextBlockId"`
	Type          string         `json:"type"`
	Props         datatypes.JSON `json:"props"`
	Content       datatypes.JSON `json:"content"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	CreatedAt     time.Time      `json:"createdAt"`
}

type GetMyBlockByIdResponseDto = BlockResponseDto

type GetMyBlocksByIdsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			BlockIds []uuid.UUID `json:"blockIds" form:"blockIds" validate:"required,min=1,max=1024"`
		},
		struct{},
	]
}

type GetMyBlocksByIdsResponseDto []BlockResponseDto

type GetMyBlocksByBlockPackIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
		},
		struct{},
	]
}

type GetMyBlocksByBlockPackIdResponseDto []BlockResponseDto
