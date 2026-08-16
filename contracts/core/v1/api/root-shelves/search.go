package apicontract

import gqlmodels "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/graphql/models"

type SearchRootShelvesRequestDto = gqlmodels.SearchRootShelfInput
type SearchRootShelvesResponseDto = gqlmodels.SearchRootShelfConnection
