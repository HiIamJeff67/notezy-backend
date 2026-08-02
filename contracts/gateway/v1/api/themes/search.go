package themesdto

import gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/graphql/models"

type SearchThemesRequestDto = gqlmodels.SearchThemeInput
type SearchThemesResponseDto = gqlmodels.SearchThemeConnection
