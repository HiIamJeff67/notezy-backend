package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/materials"

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
				apicontract.GetMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyMaterialById,
		)
		materialRoutes.POST(
			"/get-and-parent-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyMaterialAndItsParentByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyMaterialAndItsParentById,
		)
		materialRoutes.POST(
			"/get-by-parent-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyMaterialsByParentSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetMyMaterialsByParentSubShelfId,
		)
		materialRoutes.POST(
			"/get-all-by-root-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyMaterialsByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetAllMyMaterialsByRootShelfId,
		)
		materialRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateMyMaterialOperation,
			),
			authMiddleware,
			endpoint.CreateMyMaterial,
		)
		materialRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.UpdateMyMaterialById,
		)
		materialRoutes.POST(
			"/save",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SaveMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.SaveMyMaterialById,
		)
		materialRoutes.POST(
			"/move",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.MoveMyMaterialById,
		)
		materialRoutes.POST(
			"/move-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMyMaterialsByIdsOperation,
			),
			authMiddleware,
			endpoint.MoveMyMaterialsByIds,
		)
		materialRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.RestoreMyMaterialById,
		)
		materialRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyMaterialsByIdsOperation,
			),
			authMiddleware,
			endpoint.RestoreMyMaterialsByIds,
		)
		materialRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyMaterialByIdOperation,
			),
			authMiddleware,
			endpoint.DeleteMyMaterialById,
		)
		materialRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyMaterialsByIdsOperation,
			),
			authMiddleware,
			endpoint.DeleteMyMaterialsByIds,
		)
		materialRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchMaterialsOperation,
			),
			authMiddleware,
			endpoint.SearchMaterials,
		)
	}
}
