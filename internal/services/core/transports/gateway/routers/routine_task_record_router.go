package routers

import (
	"github.com/gin-gonic/gin"

	routinetaskrecordsdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/routine-task-records"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
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
				routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdOperation,
			),
			authMiddleware,
			endpoint.GetAllMyRoutineTaskRecordsByRoutineTaskId,
		)
		routineTaskRecordRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetaskrecordsdto.SearchRoutineTaskRecordsOperation,
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
				routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskRecordStatusCount,
		)
		visualizationRoutes.POST(
			"/purpose-count",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskRecordPurposeCount,
		)
		visualizationRoutes.POST(
			"/scheduled-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskRecordScheduledAtCount,
		)
		visualizationRoutes.POST(
			"/actual-started-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskRecordActualStartedAtCount,
		)
		visualizationRoutes.POST(
			"/actual-ended-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskRecordActualEndedAtCount,
		)
	}
}
