package apicontract

import (
	"github.com/google/uuid"

	gqlmodels "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/graphql/models"
)

type SearchUsersRequestDto = gqlmodels.SearchUserInput
type SearchUsersResponseDto = gqlmodels.SearchUserConnection
type LoadThemeAuthorsRequestDto []uuid.UUID
type LoadThemeAuthorsResponseDto []*gqlmodels.PublicUser
