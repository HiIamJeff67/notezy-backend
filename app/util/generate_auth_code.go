package util

import (
	"math/rand/v2"
	"strconv"

	constants "notezy-backend/app/shared/constants"
)

func GenerateAuthCode() string {
	randomNumber := rand.IntN(constants.MaxAuthCode + 1)
	stringRandomNumber := strconv.Itoa(randomNumber)
	for len(stringRandomNumber) < constants.MaxLengthOfAuthCode {
		stringRandomNumber = "0" + stringRandomNumber
	}
	return stringRandomNumber
}
