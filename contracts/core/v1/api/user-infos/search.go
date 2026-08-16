package apicontract

import (
	"github.com/google/uuid"

	gqlmodels "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/graphql/models"
)

type LoadUserInfosRequestDto []uuid.UUID
type LoadUserInfosResponseDto []*gqlmodels.PublicUserInfo
