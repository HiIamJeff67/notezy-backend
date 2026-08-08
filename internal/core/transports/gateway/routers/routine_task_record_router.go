package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routine-task-records"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

func configureRoutineTaskRecordRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.RoutineTaskRecordEndpointInterface,
) {
	routineTaskRecordRoutes := router.Group("/routine-task-records")
	{
		routineTaskRecordRoutes.POST(
			"/get-all-by-routine-task-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyRoutineTaskRecordsByRoutineTaskIdOperation,
			),
			authMiddleware,
			endpoint.GetAllMyRoutineTaskRecordsByRoutineTaskId,
		)
		routineTaskRecordRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchRoutineTaskRecordsOperation,
			),
			authMiddleware,
			endpoint.SearchRoutineTaskRecords,
		)
	}
	visualizationRoutes := routineTaskRecordRoutes.Group("/visualizations")
	{
		visualizationRoutes.POST(
			"/status-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskRecordStatusCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskRecordStatusCount,
		)
		visualizationRoutes.POST(
			"/purpose-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskRecordPurposeCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskRecordPurposeCount,
		)
		visualizationRoutes.POST(
			"/scheduled-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskRecordScheduledAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskRecordScheduledAtCount,
		)
		visualizationRoutes.POST(
			"/actual-started-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskRecordActualStartedAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskRecordActualStartedAtCount,
		)
		visualizationRoutes.POST(
			"/actual-ended-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskRecordActualEndedAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskRecordActualEndedAtCount,
		)
	}
}
