package userinfosdto

import (
	"github.com/google/uuid"

	gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/graphql/models"
)

type LoadUserInfosRequestDto []uuid.UUID
type LoadUserInfosResponseDto []*gqlmodels.PublicUserInfo
