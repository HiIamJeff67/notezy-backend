package postgres

import "testing"

func TestConnectionString(t *testing.T) {
	config := Config{
		Host:     "database",
		Port:     "5432",
		User:     "notegic",
		Name:     "notegic",
		Password: "password",
	}

	want := "host=database port=5432 user=notegic dbname=notegic password=password sslmode=disable"
	if got := ConnectionString(config); got != want {
		t.Fatalf("ConnectionString() = %q, want %q", got, want)
	}
}
