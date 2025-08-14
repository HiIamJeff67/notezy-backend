package dtos

import (
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type GetRecentShelvesReqDto struct {
	OwnerId uuid.UUID // extracted from the access token of AuthMiddleware()
	Limit   int       `json:"limit" validate:"min=1,max=100"`
}

/* ============================== Response DTO ============================== */

type GetRecentShelvesResDto struct {
	Name             string    `json:"name"`
	EncodedStructure []byte    `json:"encodedStructure"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CreatedAt        time.Time `json:"createdAt"`
}
