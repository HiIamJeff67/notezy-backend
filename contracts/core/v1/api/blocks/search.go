package blocksdto

import gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/graphql/models"

type SearchBlocksRequestDto = gqlmodels.SearchBlockInput
type SearchBlocksResponseDto = gqlmodels.SearchBlockConnection
