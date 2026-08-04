package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	themesdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/themes"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
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
	request := &gatewaycontract.Request[themesdto.SearchThemesRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.themeService.SearchPublicThemes(ctx.Request.Context(), request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[themesdto.SearchThemesResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
