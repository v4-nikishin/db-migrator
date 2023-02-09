package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/v4-nikishin/db-migrator/internal/app"
	"github.com/v4-nikishin/db-migrator/internal/config"
	"github.com/v4-nikishin/db-migrator/internal/logger"
	"github.com/v4-nikishin/db-migrator/internal/storage"
	"github.com/v4-nikishin/db-migrator/internal/version"
)

var configFile string
var sqlFile string
var migrationType string

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "Path to configuration file")
	flag.StringVar(&sqlFile, "sql", "", "Path to migration sql file")
	flag.StringVar(&migrationType, "migration", "", "Migration type (up/down)")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()
		return
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Printf("failed to configure service %s\n", err)
		os.Exit(1)
	}

	logg := logger.New(cfg.Logger, os.Stdout)
	repo, err := storage.New(context.Background(), cfg.DB, logg)
	if err != nil {
		logg.Error("failed to create sql storage: " + err.Error())
		return
	}
	defer repo.Close()

	a := app.New(logg, repo)

	logg.Info("Start migration " + migrationType)

	if err = a.Migration(sqlFile, app.MigrationType(migrationType)); err != nil {
		logg.Error(err.Error())
		return
	}

	logg.Info("Success migration " + migrationType)
}
