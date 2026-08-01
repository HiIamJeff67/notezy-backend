package usersdto

import (
	"github.com/google/uuid"

	gqlmodels "github.com/HiIamJeff67/notezy-backend/internal/platform/graphql/models"
)

type SearchUsersRequestDto = gqlmodels.SearchUserInput
type SearchUsersResponseDto = gqlmodels.SearchUserConnection
type LoadThemeAuthorsRequestDto []uuid.UUID
type LoadThemeAuthorsResponseDto []*gqlmodels.PublicUser
