package repository

import (
	"database/sql"
	"testing"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/stretchr/testify/require"
)

func TestNewDatabase(t *testing.T) {
	t.Parallel()

	log := logger.CreateLogger("Info").GetLogger()

	tests := []struct {
		name    string
		opener  SQLDBOpener
		builder DataSourceBuilder
		driver  string
		dsn     string
		wantErr error
	}{
		{
			"test 1",
			sql.Open,
			PGDataSourceBuilder,
			"pgx",
			"postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable",
			nil,
		},
		// TODO add error test cases
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewDatabase(test.opener, test.builder, log, test.driver, test.dsn)
			require.ErrorIs(t, err, test.wantErr)
		})
	}
}
