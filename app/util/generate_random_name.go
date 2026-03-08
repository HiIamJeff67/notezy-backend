package util

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
)

func GenerateRandomFakeDisplayName() string {
	gofakeit.Seed(0)
	animal := gofakeit.LastName()
	adjective := gofakeit.AdjectiveDescriptive()
	number := gofakeit.Number(100000, 999999)
	return fmt.Sprintf("%s%s%d", adjective, animal, number)
}
