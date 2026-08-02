package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	rootshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/root-shelves"
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
	request := &core.Request[rootshelvesdto.GetMyRootShelfByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.GetMyRootShelfById(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.GetMyRootShelfByIdResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) CreateRootShelf(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.CreateRootShelfRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.CreateRootShelf(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.CreateRootShelfResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) CreateRootShelves(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.CreateRootShelvesRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.CreateRootShelves(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.CreateRootShelvesResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) UpdateMyRootShelfById(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.UpdateMyRootShelfByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.UpdateMyRootShelfById(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.UpdateMyRootShelfByIdResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) UpdateMyRootShelvesByIds(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.UpdateMyRootShelvesByIds(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.UpdateMyRootShelvesByIdsResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) RestoreMyRootShelfById(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.RestoreMyRootShelfByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.RestoreMyRootShelfById(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.RestoreMyRootShelfByIdResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) RestoreMyRootShelvesByIds(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.RestoreMyRootShelvesByIds(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.RestoreMyRootShelvesByIdsResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) DeleteMyRootShelfById(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.DeleteMyRootShelfByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.DeleteMyRootShelfById(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.DeleteMyRootShelfByIdResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) DeleteMyRootShelvesByIds(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.DeleteMyRootShelvesByIds(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.DeleteMyRootShelvesByIdsResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) GetMyRootShelfPermission(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.GetMyRootShelfPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.GetMyRootShelfPermission(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.GetMyRootShelfPermissionResponseDto]{
		Version:  core.Version,
		Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RootShelfEndpoint) CreateMyRootShelfPermission(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.CreateMyRootShelfPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.CreateMyRootShelfPermission(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.CreateMyRootShelfPermissionResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RootShelfEndpoint) UpsertMyRootShelfPermission(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.UpsertMyRootShelfPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.UpsertMyRootShelfPermission(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.UpsertMyRootShelfPermissionResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RootShelfEndpoint) UpsertMyRootShelfPermissions(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.UpsertMyRootShelfPermissions(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.UpsertMyRootShelfPermissionsResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) UpdateMyRootShelfPermission(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.UpdateMyRootShelfPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.UpdateMyRootShelfPermission(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.UpdateMyRootShelfPermissionResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RootShelfEndpoint) TransferMyRootShelfOwnership(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.TransferMyRootShelfOwnershipRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.TransferMyRootShelfOwnership(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.TransferMyRootShelfOwnershipResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) DeleteMyRootShelfPermission(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.DeleteMyRootShelfPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.DeleteMyRootShelfPermission(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.DeleteMyRootShelfPermissionResponseDto]{
		Version:  core.Version,
		Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RootShelfEndpoint) DeleteMyRootShelfPermissions(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.rootShelfService.DeleteMyRootShelfPermissions(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.DeleteMyRootShelfPermissionsResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RootShelfEndpoint) LeaveMyRootShelf(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.LeaveMyRootShelfRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if exception := t.rootShelfService.LeaveMyRootShelf(ctx.Request.Context(), &request.Dto); exception != nil {
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.LeaveMyRootShelfResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: rootshelvesdto.LeaveMyRootShelfResponseDto{},
	})
}

func (t *RootShelfEndpoint) LeaveMyRootShelves(ctx *gin.Context) {
	request := &core.Request[rootshelvesdto.LeaveMyRootShelvesRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if exception := t.rootShelfService.LeaveMyRootShelves(ctx.Request.Context(), &request.Dto); exception != nil {
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

	ctx.JSON(http.StatusOK, core.Response[rootshelvesdto.LeaveMyRootShelvesResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: rootshelvesdto.LeaveMyRootShelvesResponseDto{},
	})
}
