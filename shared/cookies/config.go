package cookies

import (
	"net/http"
	"time"
)

type Config struct {
	Name     ValidCookieName
	Path     string
	Duration time.Duration
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}
