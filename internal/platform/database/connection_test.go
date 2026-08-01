package database

import (
	"testing"

	configs "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
)

func TestConnectionString(t *testing.T) {
	config := configs.DatabaseConfig{
		Host:     "database",
		Port:     "5432",
		User:     "notezy",
		DBName:   "notezy",
		Password: "password",
	}

	want := "host=database port=5432 user=notezy dbname=notezy password=password sslmode=disable"
	if got := ConnectionString(config); got != want {
		t.Fatalf("ConnectionString() = %q, want %q", got, want)
	}
}
