package cookies

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

type CookieHandlerInterface interface {
	Get(ctx *gin.Context) (string, *exceptions.Exception)
	Set(ctx *gin.Context, value string)
	Delete(ctx *gin.Context)
}

type CookieHandler struct {
	name     types.ValidCookieName
	path     string
	duration time.Duration
	secure   bool
	httpOnly bool
	sameSite http.SameSite
}

// a constructor of the cookie handler
func NewCookieHandler(name types.ValidCookieName, path string, duration time.Duration, secure, httpOnly bool, sameSite http.SameSite) *CookieHandler {
	return &CookieHandler{
		name:     name,
		path:     path,
		duration: duration,
		secure:   secure,
		httpOnly: httpOnly,
		sameSite: sameSite,
	}
}

func (h *CookieHandler) Get(ctx *gin.Context) (string, *exceptions.Exception) {
	value, err := ctx.Cookie(h.name.String())
	if err != nil {
		return "", exceptions.Cookie.NotFound(string(h.name)).WithError(err)
	}
	return value, nil
}

func (h *CookieHandler) Set(ctx *gin.Context, value string) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     h.name.String(),
		Path:     h.path,
		Expires:  time.Now().Add(h.duration),
		Secure:   h.secure,
		HttpOnly: h.httpOnly,
		SameSite: h.sameSite,
		Value:    value,
		Domain:   "",
	})
}

func (h *CookieHandler) Delete(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     h.name.String(),
		Path:     h.path,
		Expires:  time.Unix(0, 0), // set to before
		MaxAge:   -1,              // set to before
		Secure:   h.secure,
		HttpOnly: h.httpOnly,
		SameSite: h.sameSite,
	})
}
