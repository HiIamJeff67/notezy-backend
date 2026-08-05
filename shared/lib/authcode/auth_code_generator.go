package authcode

import (
	"math/rand/v2"
	"strconv"
	"time"
)

const (
	authCodeExpiration = 3 * time.Minute
	authCodeLength     = 6
	authCodeMaxValue   = 999999
)

type AuthCodeGenerator struct {
	expiration time.Duration
	length     int
	maxValue   int
}

func New() *AuthCodeGenerator {
	return &AuthCodeGenerator{
		expiration: authCodeExpiration,
		length:     authCodeLength,
		maxValue:   authCodeMaxValue,
	}
}

func (g *AuthCodeGenerator) Generate() string {
	code := strconv.Itoa(rand.IntN(g.maxValue + 1))
	for len(code) < g.length {
		code = "0" + code
	}
	return code
}

func (g *AuthCodeGenerator) ExpireAt(now time.Time) time.Time {
	return now.Add(g.expiration)
}
