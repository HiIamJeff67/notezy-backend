package routers

import (
	"github.com/gin-gonic/gin"

	materialsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/materials"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

func configureMaterialRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.MaterialEndpointInterface,
) {
	materialRoutes := router.Group("/materials")
	{
		materialRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.GetMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyMaterialById,
		)
		materialRoutes.POST(
			"/get-and-parent-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.GetMyMaterialAndItsParentByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyMaterialAndItsParentById,
		)
		materialRoutes.POST(
			"/get-by-parent-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.GetMyMaterialsByParentSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetMyMaterialsByParentSubShelfId,
		)
		materialRoutes.POST(
			"/get-all-by-root-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.GetAllMyMaterialsByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetAllMyMaterialsByRootShelfId,
		)
		materialRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.CreateMyMaterialOperation,
			),
			authMiddleware,
			endpoint.CreateMyMaterial,
		)
		materialRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.UpdateMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.UpdateMyMaterialById,
		)
		materialRoutes.POST(
			"/save",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.SaveMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.SaveMyMaterialById,
		)
		materialRoutes.POST(
			"/move",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.MoveMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.MoveMyMaterialById,
		)
		materialRoutes.POST(
			"/move-many",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.MoveMyMaterialsByIdsOperation,
			),
			authMiddleware,
			endpoint.MoveMyMaterialsByIds,
		)
		materialRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.RestoreMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.RestoreMyMaterialById,
		)
		materialRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.RestoreMyMaterialsByIdsOperation,
			),
			authMiddleware,
			endpoint.RestoreMyMaterialsByIds,
		)
		materialRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.DeleteMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.DeleteMyMaterialById,
		)
		materialRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.DeleteMyMaterialsByIdsOperation,
			),
			authMiddleware,
			endpoint.DeleteMyMaterialsByIds,
		)
		materialRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				materialsdto.SearchMaterialsOperation,
			),
			authMiddleware,
			endpoint.SearchMaterials,
		)
	}
}
