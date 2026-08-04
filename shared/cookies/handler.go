package cookies

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CookieHandlerInterface interface {
	Get(ctx *gin.Context) (string, error)
	Set(ctx *gin.Context, value string)
	Delete(ctx *gin.Context)
}

type CookieHandler struct {
	name     ValidCookieName
	path     string
	duration time.Duration
	secure   bool
	httpOnly bool
	sameSite http.SameSite
}

// a constructor of the cookie handler
func New(config Config) *CookieHandler {
	return &CookieHandler{
		name:     config.Name,
		path:     config.Path,
		duration: config.Duration,
		secure:   config.Secure,
		httpOnly: config.HTTPOnly,
		sameSite: config.SameSite,
	}
}

func (h *CookieHandler) Get(ctx *gin.Context) (string, error) {
	return ctx.Cookie(h.name.String())
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
