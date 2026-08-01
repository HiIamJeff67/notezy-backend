package routinesdto

import gqlmodels "github.com/HiIamJeff67/notezy-backend/internal/platform/graphql/models"

type SearchRoutinesRequestDto = gqlmodels.SearchRoutineInput
type SearchRoutinesResponseDto = gqlmodels.SearchRoutineConnection
