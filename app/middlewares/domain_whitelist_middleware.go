package middlewares

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	exceptions "notezy-backend/app/exceptions"
	util "notezy-backend/app/util"
)

func isAllowedOrigin(origin string, allowedDomains []string) bool {
	parsedURL, err := url.Parse(origin)
	if err != nil {
		return false
	}

	host := parsedURL.Host
	for _, allowed := range allowedDomains {
		if host == allowed {
			return true
		}
	}
	return false
}

func isAllowedReferer(referer string, allowedDomains []string) bool {
	parsedURL, err := url.Parse(referer)
	if err != nil {
		return false
	}

	host := parsedURL.Host
	for _, allowed := range allowedDomains {
		if host == allowed {
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
				ctx.AbortWithStatusJSON(http.StatusForbidden,
					exceptions.Auth.PermissionDeniedDueToInvalidRequestOriginDomain(origin).GetGinH())
				return
			}
		}

		referer := ctx.GetHeader("Referer")
		if referer != "" && origin == "" {
			if !isAllowedReferer(referer, allowedDomains) {
				ctx.AbortWithStatusJSON(http.StatusForbidden,
					exceptions.Auth.PermissionDeniedDueToInvalidRequestOriginDomain(referer).GetGinH())
				return
			}
		}

		ctx.Next()
	}
}
