package app

import (
	"fmt"
	"path"
	"time"

	"github.com/v4-nikishin/db-migrator/internal/logger"
	"github.com/v4-nikishin/db-migrator/internal/parser"
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
	switch mt {
	case MigrationUp:
		m, err := a.repo.GetMigration(name)
		if err == nil {
			return fmt.Errorf("migration already is %s", m.Status)
		}
		m = storage.Migration{
			Name:   name,
			Date:   time.Now().String(),
			Status: storage.Running,
		}
		if err = a.repo.CreateMigration(m); err != nil {
			return fmt.Errorf("failed to create %s: %w", name, err)
		}
		if err = a.migrationUp(filePath); err != nil {
			return fmt.Errorf("failed to migrate up %s: %w", name, err)
		}
		m.Status = storage.Done
		if err = a.repo.UpdateMigration(m); err != nil {
			return fmt.Errorf("failed to update %s: %w", name, err)
		}
	case MigrationDown:
		m, err := a.repo.GetMigration(name)
		if err != nil {
			return fmt.Errorf("failed to get migration %s", name)
		}
		if err := a.migrationDown(filePath); err != nil {
			m.Status = storage.Failed
			if err := a.repo.UpdateMigration(m); err != nil {
				return fmt.Errorf("failed to update %s: %w", name, err)
			}
			return fmt.Errorf("failed to migrate down %s: %w", name, err)
		}
		if err := a.repo.DeleteMigration(name); err != nil {
			return fmt.Errorf("failed to delete %s: %w", name, err)
		}
	default:
		return fmt.Errorf("unknown migration type %s ", mt)
	}
	return nil
}

func (a *App) migrationUp(filePath string) error {
	p := parser.New(a.log, filePath)
	sts, err := p.UpMigration()
	if err != nil {
		return err
	}
	for _, st := range sts {
		if _, err := a.repo.DB.Exec(st); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) migrationDown(filePath string) error {
	p := parser.New(a.log, filePath)
	sts, err := p.DownMigration()
	if err != nil {
		return err
	}
	for _, st := range sts {
		if _, err := a.repo.DB.Exec(st); err != nil {
			return err
		}
	}
	return nil
}
