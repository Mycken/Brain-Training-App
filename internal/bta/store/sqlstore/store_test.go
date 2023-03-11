package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseUrl string
)

func TestMain(m *testing.M) {
	databaseUrl = os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "host=localhost port=5432 user=postgres password=123 dbname=bta_test sslmode=disable"
	}
	os.Exit(m.Run())
}
