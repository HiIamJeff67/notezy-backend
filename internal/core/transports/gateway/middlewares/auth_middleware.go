package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"

	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/core/contexts"
	data "github.com/HiIamJeff67/notezy-backend/internal/core/data/database"
	options "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/repositories"
)

func AuthMiddleware(userRepository repositories.UserRepositoryInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userPublicId, exception := contexts.GetActorUserPublicId(ctx.Request.Context())
		if exception != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"InvalidDelegation",
					"Core",
					"AuthenticateRequest",
					"a user delegation subject is required",
					http.StatusUnauthorized,
				),
			})
			return
		}

		accessToken := ""
		accessTokenExists := false
		if cookie, err := ctx.Request.Cookie(cookies.ValidCookieName_AccessToken.String()); err == nil && strings.TrimSpace(cookie.Value) != "" {
			accessToken = cookie.Value
			accessTokenExists = true
		}

		refreshToken := ""
		refreshTokenExists := false
		if cookie, err := ctx.Request.Cookie(cookies.ValidCookieName_RefreshToken.String()); err == nil && strings.TrimSpace(cookie.Value) != "" {
			refreshToken = cookie.Value
			refreshTokenExists = true
		}
		userAgent := ctx.GetHeader("User-Agent")
		if accessTokenExists {
			claims, err := sharedtokens.ParseAccessToken(accessToken)
			if err == nil && claims.Subject == userPublicId.String() && claims.UserAgent == userAgent {
				if setActorUserId(ctx, userRepository, userPublicId) {
					ctx.Next()
				}
				return
			}
		}

		if !refreshTokenExists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"Unauthorized",
					"Core",
					"AuthenticateRequest",
					"the forwarded authentication credentials are invalid",
					http.StatusUnauthorized,
				),
			})
			return
		}

		claims, err := sharedtokens.ParseRefreshToken(refreshToken)
		if err != nil || claims.Subject != userPublicId.String() || claims.UserAgent != userAgent {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"Unauthorized",
					"Core",
					"AuthenticateRequest",
					"the forwarded authentication credentials are invalid",
					http.StatusUnauthorized,
				),
			})
			return
		}

		if userRepository == nil {
			ctx.Next()
			return
		}

		user, exception := userRepository.GetOneByPublicId(
			userPublicId,
			nil,
			options.WithDB(data.NotezyDB),
		)
		if exception != nil || user.RefreshToken != refreshToken || user.UserAgent != userAgent {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"Unauthorized",
					"Core",
					"AuthenticateRequest",
					"the forwarded refresh credential could not be authenticated",
					http.StatusUnauthorized,
				),
			})
			return
		}

		newAccessToken, err := sharedtokens.GenerateAccessToken(
			user.PublicId.String(),
			sharedtokens.AccessTokenClaims{
				Name:      user.Name,
				Email:     user.Email,
				UserAgent: user.UserAgent,
			},
		)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"GenerationFailed",
					"Core",
					"AuthenticateRequest",
					"failed to generate a new access token",
					http.StatusInternalServerError,
					true,
				).WithOrigin(err),
			})
			return
		}
		newCSRFToken, err := sharedtokens.GenerateCSRFToken(sharedtokens.CSRFTokenClaims{})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"GenerationFailed",
					"Core",
					"AuthenticateRequest",
					"failed to generate a new CSRF token",
					http.StatusInternalServerError,
					true,
				).WithOrigin(err),
			})
			return
		}

		ctx.Header(gatewaycontract.CoreAuthRefreshed.String(), "true")
		ctx.Header(gatewaycontract.CoreSetAccessToken.String(), *newAccessToken)
		ctx.Header(gatewaycontract.CoreSetCSRFToken.String(), *newCSRFToken)
		if !setActorUserId(ctx, userRepository, userPublicId) {
			return
		}
		ctx.Next()
	}
}

func setActorUserId(
	ctx *gin.Context,
	userRepository repositories.UserRepositoryInterface,
	userPublicId uuid.UUID,
) bool {
	if userRepository == nil {
		return true
	}

	user, exception := userRepository.GetOneByPublicId(
		userPublicId,
		nil,
		options.WithDB(data.NotezyDB),
	)
	if exception != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   ctx.GetHeader("X-Request-Id"),
				RespondedAt: time.Now(),
			},
			Data: struct{}{},
			Exception: exceptions.New(
				"Unauthorized",
				"Core",
				"AuthenticateRequest",
				"the delegated user subject could not be authenticated",
				http.StatusUnauthorized,
			),
		})
		return false
	}

	requestContext := contexts.WithActorUserId(ctx.Request.Context(), user.Id)
	ctx.Request = ctx.Request.WithContext(contexts.WithActorUserName(requestContext, user.Name))
	ctx.Set(sharedcontexts.ContextFieldName_User_Name.String(), user.Name)
	ctx.Set(sharedcontexts.ContextFieldName_User_Email.String(), user.Email)
	ctx.Set(sharedcontexts.ContextFieldName_User_Role.String(), user.Role)
	ctx.Set(sharedcontexts.ContextFieldName_User_Plan.String(), user.Plan)
	return true
}
