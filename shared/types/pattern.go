package types

import "github.com/gin-gonic/gin"

type ControllerFunc[DtoType any] func(ctx *gin.Context, reqDto DtoType)

type BinderFunc[DtoType any] func(ctx *gin.Context, controllerFunc ControllerFunc[DtoType])
