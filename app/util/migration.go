package util

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
)

func GenerateMigrationFileName(dbName string) string {
	return fmt.Sprintf("%s_%s_%s_%s", dbName, randomdata.City(), randomdata.Noun(), randomdata.Locale())
}