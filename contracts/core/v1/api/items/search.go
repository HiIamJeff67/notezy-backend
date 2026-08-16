package apicontract

import gqlmodels "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/graphql/models"

type SearchItemsRequestDto = gqlmodels.SearchItemInput
type SearchItemsResponseDto = gqlmodels.SearchItemConnection
