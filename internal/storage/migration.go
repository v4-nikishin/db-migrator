package storage

type MigrationStatus string

const (
	Running MigrationStatus = "running"
	Done    MigrationStatus = "done"
	Failed  MigrationStatus = "failed"
)

type Migration struct {
	Name   string
	Date   string
	Status MigrationStatus
}
