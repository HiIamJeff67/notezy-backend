package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
)

func isAllowedOrigin(origin string, allowedDomains []string) bool {
	for _, allowed := range allowedDomains {
		if origin == allowed || (origin[len(origin)-1] == '/' && origin[0:len(origin)-1] == allowed) {
			return true
		}
	}
	return false
}

func isAllowedReferer(referer string, allowedDomains []string) bool {
	for _, allowed := range allowedDomains {
		if referer == allowed || (referer[len(referer)-1] == '/' && referer[0:len(referer)-1] == allowed) {
			return true
		}
	}
	return false
}

func DomainWhiteListMiddleware() gin.HandlerFunc {
	var allowedDomains []string
	if envDomains := os.Getenv("ALLOWED_DOMAINS"); len(strings.ReplaceAll(envDomains, " ", "")) > 0 {
		additionalDomains := strings.Split(envDomains, ",")
		for _, domain := range additionalDomains {
			allowedDomains = append(allowedDomains, strings.TrimSpace(domain))
		}
	}
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if origin != "" {
			if !isAllowedOrigin(origin, allowedDomains) {
				logs.NotezyLogger.Alert(ctx.Request.Context(), nil, fmt.Sprintf("Blocked Origin: %s, allowed origins: ", origin))
				for _, domain := range allowedDomains {
					logs.NotezyLogger.Alert(ctx.Request.Context(), nil, domain)
				}
				ctx.AbortWithStatusJSON(http.StatusForbidden,
					responsewriter.GetGinH(exceptions.New(
						"PermissionDeniedDueToInvalidRequestOriginDomain",
						"Auth",
						"Authorize",
						fmt.Sprintf("The current request origin domain of %s is invalid", origin),
						http.StatusForbidden,
					)))
				return
			}
		}

		referer := ctx.GetHeader("Referer")
		if referer != "" && origin == "" {
			if !isAllowedReferer(referer, allowedDomains) {
				logs.NotezyLogger.Alert(ctx.Request.Context(), nil, fmt.Sprintf("Blocked Referer: %s", referer))
				ctx.AbortWithStatusJSON(http.StatusForbidden,
					responsewriter.GetGinH(exceptions.New(
						"PermissionDeniedDueToInvalidRequestOriginDomain",
						"Auth",
						"Authorize",
						fmt.Sprintf("The current request origin domain of %s is invalid", referer),
						http.StatusForbidden,
					)))
				return
			}
		}

		ctx.Next()
	}
}
