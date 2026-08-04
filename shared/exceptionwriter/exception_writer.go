package exceptionwriter

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/traces"
)

func IncrementMeter(exception *exceptions.Exception, ctx *gin.Context, names ...string) {
	if metrics.NotezyMeter == nil {
		return
	}

	isTotalCounted := false
	for _, name := range names {
		if name == "server.responses.failed.total" {
			isTotalCounted = true
		}
		metrics.NotezyMeter.Count(ctx, name, 1)
	}
	if !isTotalCounted {
		metrics.NotezyMeter.Count(ctx, "server.responses.failed.total", 1)
	}
}

func GetGinH(exception *exceptions.Exception) gin.H {
	exception = exception.ToPublic()

	return gin.H{
		"reason":    exception.Reason,
		"domain":    exception.Domain,
		"operation": exception.Operation,
		"message":   exception.Message,
		"retryable": exception.Retryable,
	}
}

func GetGinHBytes(exception *exceptions.Exception) ([]byte, error) {
	return json.Marshal(GetGinH(exception))
}

func GetResponseJSONBytes(exception *exceptions.Exception) ([]byte, error) {
	return json.Marshal(gin.H{
		"success":   false,
		"data":      nil,
		"exception": GetGinH(exception),
	})
}

func getRequestContext(ctx *gin.Context) context.Context {
	if ctx != nil && ctx.Request != nil {
		return ctx.Request.Context()
	}

	return context.Background()
}

func ResponseWithJSON(exception *exceptions.Exception, ctx *gin.Context, names ...string) {
	publicException := ToPublic(getRequestContext(ctx), exception)
	IncrementMeter(publicException, ctx, names...)
	if exception != nil &&
		exception.HTTPStatusCode() >= http.StatusInternalServerError &&
		traces.NotezyTracer != nil {
		traces.NotezyTracer.RecordError(ctx, exception)
	}

	ctx.JSON(publicException.HTTPStatusCode(), gin.H{
		"success":   false,
		"data":      nil,
		"exception": GetGinH(publicException),
	})
}

func SafelyResponseWithJSON(exception *exceptions.Exception, ctx *gin.Context, names ...string) {
	publicException := ToPublic(getRequestContext(ctx), exception)
	IncrementMeter(publicException, ctx, names...)
	if exception != nil &&
		exception.HTTPStatusCode() >= http.StatusInternalServerError &&
		traces.NotezyTracer != nil {
		traces.NotezyTracer.RecordError(ctx, exception)
	}

	ctx.JSON(publicException.HTTPStatusCode(), gin.H{
		"success":   false,
		"data":      nil,
		"exception": GetGinH(publicException),
	})
}

func SafelyAbortAndResponseWithJSON(exception *exceptions.Exception, ctx *gin.Context, names ...string) {
	publicException := ToPublic(getRequestContext(ctx), exception)
	IncrementMeter(publicException, ctx, names...)
	if exception != nil &&
		exception.HTTPStatusCode() >= http.StatusInternalServerError &&
		traces.NotezyTracer != nil {
		traces.NotezyTracer.RecordError(ctx, exception)
	}

	ctx.AbortWithStatusJSON(publicException.HTTPStatusCode(), gin.H{
		"success":   false,
		"data":      nil,
		"exception": GetGinH(publicException),
	})
}

func ToGraphQLError(exception *exceptions.Exception, ctx context.Context) *gqlerror.Error {
	exception = ToPublic(ctx, exception)
	extensions := map[string]any{
		"reason":     exception.Reason,
		"domain":     exception.Domain,
		"operation":  exception.Operation,
		"httpStatus": exception.HTTPStatusCode(),
		"retryable":  exception.Retryable,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	var path ast.Path
	var locations []gqlerror.Location
	if ctx != nil {
		if fieldContext := graphql.GetFieldContext(ctx); fieldContext != nil {
			path = fieldContext.Path()
			if fieldContext.Field.Position != nil {
				locations = []gqlerror.Location{{
					Line:   fieldContext.Field.Position.Line,
					Column: fieldContext.Field.Position.Column,
				}}
			}
		}
		if operationContext := graphql.GetOperationContext(ctx); operationContext != nil && operationContext.OperationName != "" {
			extensions["operationName"] = operationContext.OperationName
		}
	}

	gqlError := &gqlerror.Error{
		Message:    exception.Message,
		Path:       path,
		Locations:  locations,
		Extensions: extensions,
	}
	return gqlError
}

func ToPublic(ctx context.Context, exception *exceptions.Exception) *exceptions.Exception {
	if exception == nil {
		return exceptions.InternalServerError()
	}
	if logs.NotezyLogger != nil {
		if err := logs.NotezyLogger.JSON(
			ctx,
			slog.LevelError,
			exception.String(),
			exception,
		); err != nil {
			logs.NotezyLogger.Error(
				ctx,
				err,
				"failed to marshal exception for logging",
			)
		}
	}

	return exception.ToPublic()
}
