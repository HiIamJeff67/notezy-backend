package apitransport

import "github.com/gin-gonic/gin"

type ControllerFunc[DtoType any] func(ctx *gin.Context, reqDto DtoType)
