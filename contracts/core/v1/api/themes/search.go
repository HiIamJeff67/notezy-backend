package apicontract

import gqlmodels "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/graphql/models"

type SearchThemesRequestDto = gqlmodels.SearchThemeInput
type SearchThemesResponseDto = gqlmodels.SearchThemeConnection
