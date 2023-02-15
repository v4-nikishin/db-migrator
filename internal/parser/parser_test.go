package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/v4-nikishin/db-migrator/internal/config"
	"github.com/v4-nikishin/db-migrator/internal/logger"
)

func TestParser(t *testing.T) {
	logg := logger.New(config.LoggerConf{Level: logger.DebugStr}, os.Stdout)
	p := New(logg, "../../migrations/test_migration.sql")

	t.Run("check command Up", func(t *testing.T) {
		const expected0 = `CREATE TABLE test_migration (
    id              serial primary key,
    name            text,
    date            text
);`
		const expected1 = `INSERT INTO test_migration (name, date) 
VALUES ('test_name', '2023-03-12');`
		sts, err := p.UpMigration()
		require.NoError(t, err)
		require.Equal(t, len(sts), 2)
		require.Equal(t, expected0, sts[0])
		require.Equal(t, expected1, sts[1])
	})
	t.Run("check command Down", func(t *testing.T) {
		const expected = "DROP TABLE test_migration;"
		sts, err := p.DownMigration()
		require.NoError(t, err)
		require.Equal(t, len(sts), 1)
		require.Equal(t, expected, sts[0])
	})
}
