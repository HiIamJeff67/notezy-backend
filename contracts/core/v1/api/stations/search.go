package apicontract

import gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/graphql/models"

type SearchStationsRequestDto = gqlmodels.SearchStationInput
type SearchStationsResponseDto = gqlmodels.SearchStationConnection
