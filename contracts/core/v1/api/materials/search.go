package apicontract

import gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/graphql/models"

type SearchMaterialsRequestDto = gqlmodels.SearchMaterialInput
type SearchMaterialsResponseDto = gqlmodels.SearchMaterialConnection
