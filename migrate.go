// Package migrate provides database migration functionality.
// It is a fork of golang-migrate/migrate with additional features and fixes.
package migrate

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// DefaultPrefetchMigrations is the default number of migrations to prefetch.
const DefaultPrefetchMigrations = 10

// ErrNoChange is returned when no migration is needed.
var ErrNoChange = errors.New("no change")

// ErrNilVersion is returned when the version is nil.
var ErrNilVersion = errors.New("nil version")

// ErrLocked is returned when the database is locked.
var ErrLocked = errors.New("database locked")

// ErrNotLocked is returned when trying to unlock a database that is not locked.
var ErrNotLocked = errors.New("database not locked")

// Migrate is the main struct for managing database migrations.
type Migrate struct {
	// sourceName is the registered source driver name.
	sourceName string
	// sourceDrv is the source driver instance.
	sourceDrv Source

	// databaseName is the registered database driver name.
	databaseName string
	// databaseDrv is the database driver instance.
	databaseDrv Database

	// Log is an optional logger.
	Log Logger

	// GracefulStop is a channel to signal a graceful stop.
	GracefulStop chan bool
	isGracefulStop bool

	isLockedMu sync.Mutex
	isLocked   bool

	// PrefetchMigrations is the number of migrations to prefetch.
	PrefetchMigrations uint
}

// Logger is the interface for logging migration activity.
type Logger interface {
	Printf(format string, v ...interface{})
	Verbose() bool
}

// New returns a new Migrate instance from the provided source and database URLs.
func New(sourceURL, databaseURL string) (*Migrate, error) {
	m := &Migrate{
		GracefulStop:       make(chan bool, 1),
		PrefetchMigrations: DefaultPrefetchMigrations,
	}

	sourceDrv, err := newSource(sourceURL, m)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	m.sourceDrv = sourceDrv

	databaseDrv, err := newDatabase(databaseURL, m)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	m.databaseDrv = databaseDrv

	return m, nil
}

// Close closes the source and database drivers.
func (m *Migrate) Close() (sourceErr error, databaseErr error) {
	if m.sourceDrv != nil {
		sourceErr = m.sourceDrv.Close()
	}
	if m.databaseDrv != nil {
		databaseErr = m.databaseDrv.Close()
	}
	return
}

// log prints a message if a logger is set.
func (m *Migrate) log(format string, v ...interface{}) {
	if m.Log != nil {
		m.Log.Printf(format, v...)
	}
}

// logVerbose prints a verbose message if a logger is set and verbose mode is enabled.
func (m *Migrate) logVerbose(format string, v ...interface{}) {
	if m.Log != nil && m.Log.Verbose() {
		m.Log.Printf(format, v...)
	}
}

// isGracefulStopSignalled checks if a graceful stop has been signalled.
func (m *Migrate) isGracefulStopSignalled() bool {
	select {
	case <-m.GracefulStop:
		return true
	default:
		return false
	}
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
