package materialsdto

import gqlmodels "github.com/HiIamJeff67/notezy-backend/internal/platform/graphql/models"

type SearchMaterialsRequestDto = gqlmodels.SearchMaterialInput
type SearchMaterialsResponseDto = gqlmodels.SearchMaterialConnection
