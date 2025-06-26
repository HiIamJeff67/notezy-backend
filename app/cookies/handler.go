package cookies

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	exceptions "notezy-backend/app/exceptions"
	shared "notezy-backend/shared"
)

type CookieHandlerInterface interface {
	GetCookie(ctx *gin.Context) (string, *exceptions.Exception)
	SetCookie(ctx *gin.Context, value string)
}

type CookieHandler struct {
	name     shared.ValidCookieName
	path     string
	expires  time.Time
	secure   bool
	httpOnly bool
	sameSite http.SameSite
}

// a constructor of the cookie handler
func NewCookieHandler(name shared.ValidCookieName, path string, expires time.Time, secure, httpOnly bool, sameSite http.SameSite) *CookieHandler {
	return &CookieHandler{
		name:     name,
		path:     path,
		expires:  expires,
		secure:   secure,
		httpOnly: httpOnly,
		sameSite: sameSite,
	}
}

func (h *CookieHandler) GetCookie(ctx *gin.Context) (string, *exceptions.Exception) {
	value, err := ctx.Cookie(h.name.String())
	if err != nil {
		return "", exceptions.Cookie.NotFound(h.name).WithError(err)
	}
	return value, nil
}

func (h *CookieHandler) SetCookie(ctx *gin.Context, value string) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     h.name.String(),
		Path:     h.path,
		Expires:  h.expires,
		Secure:   h.secure,
		HttpOnly: h.httpOnly,
		SameSite: h.sameSite,
		Value:    value,
	})
}
