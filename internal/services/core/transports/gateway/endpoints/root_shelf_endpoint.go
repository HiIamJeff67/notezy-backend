package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	rootshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/root-shelves"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type RootShelfEndpointInterface interface {
	GetMyRootShelfById(ctx *gin.Context)
	CreateRootShelf(ctx *gin.Context)
	CreateRootShelves(ctx *gin.Context)
	UpdateMyRootShelfById(ctx *gin.Context)
	UpdateMyRootShelvesByIds(ctx *gin.Context)
	RestoreMyRootShelfById(ctx *gin.Context)
	RestoreMyRootShelvesByIds(ctx *gin.Context)
	DeleteMyRootShelfById(ctx *gin.Context)
	DeleteMyRootShelvesByIds(ctx *gin.Context)
	GetMyRootShelfPermission(ctx *gin.Context)
	CreateMyRootShelfPermission(ctx *gin.Context)
	UpsertMyRootShelfPermission(ctx *gin.Context)
	UpsertMyRootShelfPermissions(ctx *gin.Context)
	UpdateMyRootShelfPermission(ctx *gin.Context)
	TransferMyRootShelfOwnership(ctx *gin.Context)
	DeleteMyRootShelfPermission(ctx *gin.Context)
	DeleteMyRootShelfPermissions(ctx *gin.Context)
	LeaveMyRootShelf(ctx *gin.Context)
	LeaveMyRootShelves(ctx *gin.Context)

	/* ============================== GraphQL Methods ============================== */
	SearchRootShelves(ctx *gin.Context)
}

type RootShelfEndpoint struct {
	rootShelfService services.RootShelfServiceInterface
}

func NewRootShelfEndpoint(
	service services.RootShelfServiceInterface,
) RootShelfEndpointInterface {
	return &RootShelfEndpoint{
		rootShelfService: service,
	}
}

/* ============================== RootShelf Endpoint Methods ============================== */

func (t *RootShelfEndpoint) GetMyRootShelfById(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.GetMyRootShelfByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.GetMyRootShelfById(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.GetMyRootShelfByIdResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) CreateRootShelf(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.CreateRootShelfRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.CreateRootShelf(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.CreateRootShelfResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) CreateRootShelves(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.CreateRootShelvesRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.CreateRootShelves(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.CreateRootShelvesResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) UpdateMyRootShelfById(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.UpdateMyRootShelfByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.UpdateMyRootShelfById(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.UpdateMyRootShelfByIdResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) UpdateMyRootShelvesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.UpdateMyRootShelvesByIds(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.UpdateMyRootShelvesByIdsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) RestoreMyRootShelfById(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.RestoreMyRootShelfByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.RestoreMyRootShelfById(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.RestoreMyRootShelfByIdResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) RestoreMyRootShelvesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.RestoreMyRootShelvesByIds(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.RestoreMyRootShelvesByIdsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) DeleteMyRootShelfById(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.DeleteMyRootShelfByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.DeleteMyRootShelfById(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.DeleteMyRootShelfByIdResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) DeleteMyRootShelvesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.DeleteMyRootShelvesByIds(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.DeleteMyRootShelvesByIdsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) GetMyRootShelfPermission(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.GetMyRootShelfPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.GetMyRootShelfPermission(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.GetMyRootShelfPermissionResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RootShelfEndpoint) CreateMyRootShelfPermission(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.CreateMyRootShelfPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.CreateMyRootShelfPermission(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.CreateMyRootShelfPermissionResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RootShelfEndpoint) UpsertMyRootShelfPermission(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.UpsertMyRootShelfPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.UpsertMyRootShelfPermission(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.UpsertMyRootShelfPermissionResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RootShelfEndpoint) UpsertMyRootShelfPermissions(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.UpsertMyRootShelfPermissions(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.UpsertMyRootShelfPermissionsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) UpdateMyRootShelfPermission(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.UpdateMyRootShelfPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.UpdateMyRootShelfPermission(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.UpdateMyRootShelfPermissionResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RootShelfEndpoint) TransferMyRootShelfOwnership(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.TransferMyRootShelfOwnershipRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.TransferMyRootShelfOwnership(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.TransferMyRootShelfOwnershipResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) DeleteMyRootShelfPermission(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.DeleteMyRootShelfPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.DeleteMyRootShelfPermission(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.DeleteMyRootShelfPermissionResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RootShelfEndpoint) DeleteMyRootShelfPermissions(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.DeleteMyRootShelfPermissions(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.DeleteMyRootShelfPermissionsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) LeaveMyRootShelf(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.LeaveMyRootShelfRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if exception := t.rootShelfService.LeaveMyRootShelf(ctx.Request.Context(), &request.Dto); exception != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.LeaveMyRootShelfResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: rootshelvesdto.LeaveMyRootShelfResponseDto{},
	})
}

func (t *RootShelfEndpoint) LeaveMyRootShelves(ctx *gin.Context) {
	request := &gatewaycontract.Request[rootshelvesdto.LeaveMyRootShelvesRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if exception := t.rootShelfService.LeaveMyRootShelves(ctx.Request.Context(), &request.Dto); exception != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[rootshelvesdto.LeaveMyRootShelvesResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: rootshelvesdto.LeaveMyRootShelvesResponseDto{},
	})
}
