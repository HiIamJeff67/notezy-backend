package themesdto

import gqlmodels "github.com/HiIamJeff67/notezy-backend/internal/platform/graphql/models"

type SearchThemesRequestDto = gqlmodels.SearchThemeInput
type SearchThemesResponseDto = gqlmodels.SearchThemeConnection
