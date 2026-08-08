package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-infos"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
)

func (t *UserInfoEndpoint) LoadUserInfos(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.LoadUserInfosRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDtos, exception := t.userInfoService.GetPublicUserInfosByUserPublicIds(ctx.Request.Context(), request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.LoadUserInfosResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: responseDtos,
	})
}
