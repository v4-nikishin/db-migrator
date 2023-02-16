package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/v4-nikishin/db-migrator/internal/app"
	"github.com/v4-nikishin/db-migrator/internal/config"
	"github.com/v4-nikishin/db-migrator/internal/logger"
	"github.com/v4-nikishin/db-migrator/internal/storage"
)

func TestMigrator(t *testing.T) {
	conninfo := "host=localhost port=5432 user=postgres password=postgres sslmode=disable"
	db, err := sql.Open("pgx", conninfo)
	require.NoError(t, err)

	dbName := "db_migrator_testing"
	_, err = db.Exec("create database " + dbName)
	defer func() {
		_, err = db.Exec("drop database " + dbName)
		require.NoError(t, err)
	}()

	logg := logger.New(config.LoggerConf{Level: logger.DebugStr}, os.Stdout)
	repo, err := storage.New(context.Background(),
		config.DBConf{DSN: fmt.Sprintf("%s dbname=%s", conninfo, dbName)}, logg)
	require.NoError(t, err)
	defer repo.Close()

	a := app.New(logg, repo)

	t.Run("check unknown migration type", func(t *testing.T) {
		err := a.Migration("/tmp/qqq", "qqq")
		require.Error(t, err)
	})
	t.Run("check migration up", func(t *testing.T) {
		err := a.Migration("../migrations/test_migration.sql", app.MigrationUp)
		require.NoError(t, err)
	})
	t.Run("check migration down", func(t *testing.T) {
		err := a.Migration("../migrations/test_migration.sql", app.MigrationDown)
		require.NoError(t, err)
	})
}
