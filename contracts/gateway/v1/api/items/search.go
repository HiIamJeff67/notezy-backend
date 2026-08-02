package itemsdto

import gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/graphql/models"

type SearchItemsRequestDto = gqlmodels.SearchItemInput
type SearchItemsResponseDto = gqlmodels.SearchItemConnection
