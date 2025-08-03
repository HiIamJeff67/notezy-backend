package util

import (
	"math/rand/v2"
	"strconv"
)

const (
	MaxLengthOfAuthCode int = 6
	MaxAuthCode         int = 999999
)

func GenerateAuthCode() string {
	randomNumber := rand.IntN(MaxAuthCode + 1)
	stringRandomNumber := strconv.Itoa(randomNumber)
	for len(stringRandomNumber) < MaxLengthOfAuthCode {
		stringRandomNumber = "0" + stringRandomNumber
	}
	return stringRandomNumber
}
