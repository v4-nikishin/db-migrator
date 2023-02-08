package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib" // To use pgx driver
	"github.com/stretchr/testify/require"
	"github.com/v4-nikishin/db-migrator/internal/config"
	"github.com/v4-nikishin/db-migrator/internal/logger"
	"github.com/v4-nikishin/db-migrator/internal/storage"
)

func TestStorage(t *testing.T) {
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

	err = repo.InitDB()
	require.NoError(t, err)

	f := "2006-01-02 15:04:05"

	now := time.Now().UTC()
	dateTime := now.Format(f)

	s := strings.Fields(dateTime)
	date := s[0]
	uuid := (uuid.New()).String()
	m := storage.Migration{
		Name:   uuid,
		Date:   date,
		Status: storage.Running,
	}

	t.Run("check create", func(t *testing.T) {
		err := repo.CreateMigration(m)
		require.NoError(t, err)
	})
	t.Run("invalid get", func(t *testing.T) {
		_, err := repo.GetMigration("QQQ")
		require.Error(t, err)
	})
	t.Run("check get", func(t *testing.T) {
		m, err := repo.GetMigration(m.Name)
		require.NoError(t, err)
		require.Equal(t, m.Name, uuid)
		require.Equal(t, m.Date, date)
		require.Equal(t, m.Status, storage.Running)
	})
	t.Run("check update", func(t *testing.T) {
		err := repo.UpdateMigration(storage.Migration{Name: uuid, Date: date, Status: storage.Done})
		require.NoError(t, err)
	})
	t.Run("check get", func(t *testing.T) {
		m, err := repo.GetMigration(m.Name)
		require.NoError(t, err)
		require.Equal(t, m.Status, storage.Done)
	})
	t.Run("check list", func(t *testing.T) {
		migrations, err := repo.Migrations()
		require.NoError(t, err)
		require.Equal(t, len(migrations), 1)
	})
	t.Run("check delete", func(t *testing.T) {
		err := repo.DeleteMigration(m.Name)
		require.NoError(t, err)
	})
	t.Run("check list", func(t *testing.T) {
		migrations, err := repo.Migrations()
		require.NoError(t, err)
		require.Equal(t, len(migrations), 0)
	})
}
