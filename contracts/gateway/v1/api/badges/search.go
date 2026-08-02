package badgesdto

import (
	"github.com/google/uuid"

	gqlmodels "github.com/HiIamJeff67/notezy-backend/contracts/graphql/models"
)

type LoadUserBadgesRequestDto []uuid.UUID
type LoadUserBadgesResponseDto []*gqlmodels.PublicBadge
