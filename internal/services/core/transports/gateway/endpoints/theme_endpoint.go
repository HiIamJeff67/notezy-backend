package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	themesdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/themes"
	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type ThemeEndpointInterface interface {
	SearchThemes(ctx *gin.Context)
}

type ThemeEndpoint struct {
	themeService services.ThemeServiceInterface
}

func NewThemeEndpoint(
	themeService services.ThemeServiceInterface,
) ThemeEndpointInterface {
	return &ThemeEndpoint{
		themeService: themeService,
	}
}

func (t *ThemeEndpoint) SearchThemes(ctx *gin.Context) {
	request := &core.Request[themesdto.SearchThemesRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.themeService.SearchPublicThemes(ctx.Request.Context(), request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[themesdto.SearchThemesResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
