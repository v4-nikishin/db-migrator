package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib" // To use pgx driver
	"github.com/v4-nikishin/db-migrator/internal/config"
	"github.com/v4-nikishin/db-migrator/internal/logger"
)

const (
	queryInitDB = `CREATE TABLE if not exists migrations (
		id              serial primary key,
		name            text,
		date            text,
		status          text
	);`
	queryCreateMigration = "insert into migrations (name, date, status) values ($1, $2, $3)"
	queryGetMigration    = "select name, date, status from migrations where name = $1"
	queryUpdateMigration = "update migrations set date=$1, status=$2 where name = $3"
	queryDeleteMigration = "delete from migrations where name = $1"
	queryMigrations      = "select name, date, status from migrations"
)

type Storage struct {
	ctx  context.Context
	cfg  config.DBConf
	logg *logger.Logger
	db   *sql.DB
}

func New(ctx context.Context, cfg config.DBConf, logger *logger.Logger) (*Storage, error) {
	s := &Storage{ctx: ctx, cfg: cfg, logg: logger}
	if err := s.connect(s.cfg.DSN); err != nil {
		return nil, fmt.Errorf("cannot connect to psql: %w", err)
	}
	return s, nil
}

func (s *Storage) connect(dsn string) (err error) {
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	return s.db.PingContext(s.ctx)
}

func (s *Storage) InitDB() error {
	_, err := s.db.ExecContext(s.ctx, queryInitDB)
	if err != nil {
		return fmt.Errorf("cannot create migrations table %w", err)
	}
	return nil
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		s.logg.Error(fmt.Sprintf("cannot close psql connection: %v", err))
	}
}

func (s *Storage) CreateMigration(m Migration) error {
	_, err := s.db.ExecContext(s.ctx, queryCreateMigration, m.Name, m.Date, m.Status)
	if err != nil {
		return fmt.Errorf("cannot add event %w", err)
	}
	return nil
}

func (s *Storage) GetMigration(name string) (Migration, error) {
	row := s.db.QueryRowContext(s.ctx, queryGetMigration, name)
	m := Migration{}
	if err := row.Scan(&m.Name, &m.Date, &m.Status); err != nil {
		return m, fmt.Errorf("cannot get migration %w", err)
	}
	return m, nil
}

func (s *Storage) UpdateMigration(m Migration) error {
	_, err := s.db.ExecContext(s.ctx, queryUpdateMigration, m.Date, m.Status, m.Name)
	if err != nil {
		return fmt.Errorf("cannot update migration %w", err)
	}
	return nil
}

func (s *Storage) DeleteMigration(name string) error {
	_, err := s.db.ExecContext(s.ctx, queryDeleteMigration, name)
	if err != nil {
		return fmt.Errorf("cannot delete migration %w", err)
	}
	return nil
}

func (s *Storage) Migrations() ([]Migration, error) {
	rows, err := s.db.QueryContext(s.ctx, queryMigrations)
	if err != nil {
		return nil, fmt.Errorf("cannot select: %w", err)
	}
	defer rows.Close()

	var migrations []Migration

	for rows.Next() {
		var m Migration
		if err := rows.Scan(
			&m.Name,
			&m.Date,
			&m.Status,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		migrations = append(migrations, m)
	}
	return migrations, rows.Err()
}
