package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	util "notezy-backend/app/util"
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

func DomainWhitelistMiddleware() gin.HandlerFunc {
	var allowedDomains []string
	if envDomains := util.GetEnv("ALLOWED_DOMAINS", ""); len(strings.ReplaceAll(envDomains, " ", "")) > 0 {
		additionalDomains := strings.Split(envDomains, ",")
		for _, domain := range additionalDomains {
			allowedDomains = append(allowedDomains, strings.TrimSpace(domain))
		}
	}
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if origin != "" {
			if !isAllowedOrigin(origin, allowedDomains) {
				logs.FAlert("Blocked Origin: %s, allowed origins: ", origin)
				for _, domain := range allowedDomains {
					logs.Alert(domain)
				}
				ctx.AbortWithStatusJSON(http.StatusForbidden,
					exceptions.Auth.PermissionDeniedDueToInvalidRequestOriginDomain(origin).GetGinH())
				return
			}
		}

		referer := ctx.GetHeader("Referer")
		if referer != "" && origin == "" {
			if !isAllowedReferer(referer, allowedDomains) {
				logs.FAlert("Blocked Referer: %s", referer)
				ctx.AbortWithStatusJSON(http.StatusForbidden,
					exceptions.Auth.PermissionDeniedDueToInvalidRequestOriginDomain(referer).GetGinH())
				return
			}
		}

		ctx.Next()
	}
}
