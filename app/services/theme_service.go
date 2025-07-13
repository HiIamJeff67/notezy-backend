package services

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	exceptions "notezy-backend/app/exceptions"
	gqlmodels "notezy-backend/app/graphql/models"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
)

/* ============================== Interface & Instance ============================== */

type ThemeServiceInterface interface {
	// services for public themes
	GetPublicThemeByEncodedSearchCursor(ctx context.Context, encodedSearchCursor string) (*gqlmodels.PublicTheme, *exceptions.Exception)
	SearchPublicThemes(ctx context.Context, gqlInput gqlmodels.SearchableThemeInput) (*gqlmodels.SearchableThemeConnection, *exceptions.Exception)
}

type ThemeService struct {
	db *gorm.DB
}

func NewThemeService(db *gorm.DB) ThemeServiceInterface {
	return &ThemeService{
		db: db,
	}
}

/* ============================== Services for Themes ============================== */

// get the theme which are created by the current user
func (s *ThemeService) GetMyThemes() {}

/* ============================== Services for Public Themes ============================== */

func (s *ThemeService) GetPublicThemeByEncodedSearchCursor(ctx context.Context, encodedSearchCursor string) (*gqlmodels.PublicTheme, *exceptions.Exception) {
	themeRepository := repositories.NewThemeRepository(s.db)

	searchCursor, exception := util.DecodeSearchCursor[gqlmodels.SearchableThemeCursorFields](encodedSearchCursor)
	if exception != nil {
		return nil, exception
	}

	publicTheme, exception := themeRepository.GetPublicOneByEncodedSearchCursor(searchCursor.Fields.EncodedSearchCursor)
	if exception != nil {
		return nil, exception
	}

	return publicTheme, nil
}

func (s *ThemeService) SearchPublicThemes(ctx context.Context, gqlInput gqlmodels.SearchableThemeInput) (*gqlmodels.SearchableThemeConnection, *exceptions.Exception) {
	startTime := time.Now()

	query := s.db.WithContext(ctx).Model(&schemas.Theme{})

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			"name ILIKE ?",
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, exception := util.DecodeSearchCursor[gqlmodels.SearchableThemeCursorFields](*gqlInput.After)
		if exception != nil {
			return nil, exception
		}

		query.Where("encoded_search_cursor > ?", searchCursor.Fields.EncodedSearchCursor)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		cending := "ASC"
		if *gqlInput.SortOrder == gqlmodels.SearchableSortOrderDesc {
			cending = "DESC"
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchableThemeSortByName:
			query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchableThemeSortByLastUpdate:
			query.Order("updated_at " + cending).
				Order("name " + cending).
				Order("created_at " + cending)
		case gqlmodels.SearchableThemeSortByCreatedAt:
			query.Order("created_at " + cending).
				Order("name " + cending).
				Order("updated_at " + cending)
		default:
			query.Order("name " + cending).
				Order("updated_at " + cending).
				Order("created_at " + cending)
		}
	}

	limit := 10
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	query = query.Limit(limit + 1)

	var themes []schemas.Theme
	if err := query.Find(&themes).Error; err != nil {
		return nil, exceptions.User.NotFound().WithError(err)
	}

	hasNextPage := len(themes) > limit
	searchEdges := make([]*gqlmodels.SearchableThemeEdge, len(themes))

	for index, theme := range themes {
		searchCursor := util.SearchCursor[gqlmodels.SearchableThemeCursorFields]{
			Fields: gqlmodels.SearchableThemeCursorFields{
				EncodedSearchCursor: theme.EncodedSearchCursor,
			},
		}
		encodedSearchCursor, exception := searchCursor.EncodeSearchCursor()
		if exception != nil {
			return nil, exception
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Searchable.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchableThemeEdge{
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

	return &gqlmodels.SearchableThemeConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
