package apicontract

import gqlmodels "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/graphql/models"

type SearchRoutineTagsRequestDto = gqlmodels.SearchRoutineTagInput
type SearchRoutineTagsResponseDto = gqlmodels.SearchRoutineTagConnection
