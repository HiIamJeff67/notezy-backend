package stationsdto

import gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/graphql/models"

type SearchStationsRequestDto = gqlmodels.SearchStationInput
type SearchStationsResponseDto = gqlmodels.SearchStationConnection
