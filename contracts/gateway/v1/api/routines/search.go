package routinesdto

import gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/graphql/models"

type SearchRoutinesRequestDto = gqlmodels.SearchRoutineInput
type SearchRoutinesResponseDto = gqlmodels.SearchRoutineConnection
