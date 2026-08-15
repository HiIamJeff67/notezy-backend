package postgres

import "testing"

func TestConnectionString(t *testing.T) {
	config := Config{
		Host:     "database",
		Port:     "5432",
		User:     "notezy",
		Name:     "notezy",
		Password: "password",
	}

	want := "host=database port=5432 user=notezy dbname=notezy password=password sslmode=disable"
	if got := ConnectionString(config); got != want {
		t.Fatalf("ConnectionString() = %q, want %q", got, want)
	}
}
