package services

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	exceptions "notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	schemas "notezy-backend/app/models/schemas"
	constants "notezy-backend/shared/constants"
	searchcursor "notezy-backend/shared/lib/searchcursor"
)

/* ============================== Interface & Instance ============================== */

type ThemeServiceInterface interface {
	// services for public themes
	GetPublicThemeByPublicId(ctx context.Context, publicId string) (*gqlmodels.PublicTheme, *exceptions.Exception)
	SearchPublicThemes(ctx context.Context, gqlInput gqlmodels.SearchThemeInput) (*gqlmodels.SearchThemeConnection, *exceptions.Exception)
}

type ThemeService struct {
	db *gorm.DB
}

func NewThemeService(db *gorm.DB) ThemeServiceInterface {
	return &ThemeService{
		db: db,
	}
}

/* ============================== Service Methods for Themes ============================== */

// get the theme which are created by the current user
func (s *ThemeService) GetMyThemeById() {}

/* ============================== Service Methods for Public Themes ============================== */

func (s *ThemeService) GetPublicThemeByPublicId(
	ctx context.Context, publicId string,
) (*gqlmodels.PublicTheme, *exceptions.Exception) {
	db := s.db.WithContext(ctx)

	theme := schemas.Theme{}
	result := db.WithContext(ctx).
		Model(&schemas.Theme{}).
		Where("public_id = ?", publicId).
		First(&theme)
	if err := result.Error; err != nil {
		return nil, exceptions.Theme.NotFound().WithError(err)
	}

	return theme.ToPublicTheme(), nil
}

func (s *ThemeService) SearchPublicThemes(
	ctx context.Context,
	gqlInput gqlmodels.SearchThemeInput,
) (*gqlmodels.SearchThemeConnection, *exceptions.Exception) {
	startTime := time.Now()

	db := s.db.WithContext(ctx)

	query := db.Model(&schemas.Theme{})

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"name ILIKE ?",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, exception := searchcursor.Decode[gqlmodels.SearchThemeCursorFields](*gqlInput.After)
		if exception != nil {
			return nil, exception
		}

		query.Where("public_id > ?", searchCursor.Fields.PublicID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		var cending string = gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchThemeSortByName:
			query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchThemeSortByLastUpdate:
			query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchThemeSortByCreatedAt:
			query.Order("created_at " + cending).
				Order("name " + cending).
				Order("updated_at " + cending)
		default:
			query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = max(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var themes []schemas.Theme
	if err := query.Find(&themes).Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	hasNextPage := len(themes) > limit
	searchEdges := make([]*gqlmodels.SearchThemeEdge, len(themes))

	for index, theme := range themes {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchThemeCursorFields]{
			Fields: gqlmodels.SearchThemeCursorFields{
				PublicID: theme.PublicId,
			},
		}
		encodedSearchCursor, exception := searchCursor.Encode()
		if exception != nil {
			return nil, exception
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchThemeEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                theme.ToPublicTheme(),
		}
	}

	searchPageInfo := &gqlmodels.SearchPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0,
	}

	if len(searchEdges) > 0 {
		searchPageInfo.StartEncodedSearchCursor = &searchEdges[0].EncodedSearchCursor
		searchPageInfo.EndEncodedSearchCursor = &searchEdges[len(searchEdges)-1].EncodedSearchCursor
	}

	searchTime := float64(time.Since(startTime).Nanoseconds()) / 1e6
	if hasNextPage {
		searchEdges = searchEdges[:limit]
	}

	return &gqlmodels.SearchThemeConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
