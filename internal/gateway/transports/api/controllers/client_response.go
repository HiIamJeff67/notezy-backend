package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
)

func writeClientResponse[D any](ctx *gin.Context, data D) {
	ctx.JSON(http.StatusOK, gatewaycontract.ClientResponse[D]{
		Success: true,
		Data:    data,
	})
}
