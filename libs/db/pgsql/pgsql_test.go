package pgsql

import (
	"testing"

	"github.com/zainulbr/simple-loan-engine/settings"
)

func TestDefaultPostgresDB(t *testing.T) {
	option := settings.PostgresOption{
		URI:     "postgres://admin:admin@localhost:5432/loan-db?sslmode=disable",
		Enabled: true,
	}

	defer Close()

	if err := Create(&option, defaultKey); err != nil {
		t.Error(err)
	}
}
