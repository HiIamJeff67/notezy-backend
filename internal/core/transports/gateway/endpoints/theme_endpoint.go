package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/themes"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	otherservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/other"
)

type ThemeEndpointInterface interface {
	SearchThemes(ctx *gin.Context)
}

type ThemeEndpoint struct {
	themeService otherservices.ThemeServiceInterface
}

func NewThemeEndpoint(
	themeService otherservices.ThemeServiceInterface,
) ThemeEndpointInterface {
	return &ThemeEndpoint{
		themeService: themeService,
	}
}

func (t *ThemeEndpoint) SearchThemes(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.SearchThemesRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.SearchThemesResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
