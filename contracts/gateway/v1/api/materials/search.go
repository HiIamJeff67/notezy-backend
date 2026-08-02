package materialsdto

import gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/graphql/models"

type SearchMaterialsRequestDto = gqlmodels.SearchMaterialInput
type SearchMaterialsResponseDto = gqlmodels.SearchMaterialConnection
