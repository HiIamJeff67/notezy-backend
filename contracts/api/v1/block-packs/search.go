package blockpacksdto

import gqlmodels "github.com/HiIamJeff67/notezy-backend/internal/platform/graphql/models"

type SearchBlockPacksRequestDto = gqlmodels.SearchBlockPackInput
type SearchBlockPacksResponseDto = gqlmodels.SearchBlockPackConnection
