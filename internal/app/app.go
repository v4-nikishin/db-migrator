package app

import (
	"fmt"
	"path"

	"github.com/v4-nikishin/db-migrator/internal/logger"
	"github.com/v4-nikishin/db-migrator/internal/storage"
)

type MigrationType string

const (
	MigrationUp   MigrationType = "up"
	MigrationDown MigrationType = "down"
)

type App struct {
	log  *logger.Logger
	repo *storage.Storage
}

func New(logger *logger.Logger, repo *storage.Storage) *App {
	return &App{log: logger, repo: repo}
}

func (a *App) Migration(filePath string, mt MigrationType) error {
	name := path.Base(filePath)
	m, err := a.repo.GetMigration(name)
	if err == nil {
		return fmt.Errorf("migration already is %s", m.Status)
	}
	switch mt {
	case MigrationUp:
		if err := a.migrationUp(filePath); err != nil {
			return fmt.Errorf("failed to migrate up %s: %w", name, err)
		}
	case MigrationDown:
		if err := a.migrationDown(filePath); err != nil {
			return fmt.Errorf("failed to migrate down %s: %w", name, err)
		}
	default:
		return fmt.Errorf("unknown migration type %s ", mt)
	}
	return nil
}

func (a *App) migrationUp(filePath string) error {
	return fmt.Errorf("not implemented")
}

func (a *App) migrationDown(filePath string) error {
	return fmt.Errorf("not implemented")
}
