package apicontract

import gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/graphql/models"

type SearchSubShelvesRequestDto = gqlmodels.SearchSubShelfInput
type SearchSubShelvesResponseDto = gqlmodels.SearchSubShelfConnection
