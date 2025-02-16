package database

import (
	"fmt"
	"time"
)

const (
	TCP      = "tcp"      // TCP connection type
	Unix     = "unix"     // Unix socket connection type
	MySQL    = "mysql"    // MySQL driver name
	Postgres = "postgres" // PostgreSQL driver name
	SQLite3  = "sqlite3"  // SQLite3 driver name
)

// DriverFactory is a function that creates a database driver.
type DriverFactory func(driver string, dsn string) (DB, error)

// ConnectConfig holds the configuration for the database connection.
type ConnectConfig struct {
	User            string        // Database user
	Password        string        // Database password
	Host            string        // Database host
	Port            int           // Database port
	Database        string        // Database name
	SocketDirectory string        // Unix socket directory
	SocketName      string        // Unix socket name
	Parameters      string        // Connection parameters
	ConnectionType  string        // Connection type
	ConnMaxLifetime time.Duration // Connection max lifetime
	ConnMaxIdleTime time.Duration // Connection max idle time
	MaxOpenConns    int           // Max open connections
	MaxIdleConns    int           // Max idle connections
	Driver          string        // Driver name

	// DSNFormat is an optional format string (e.g. "%s:%s@tcp(%s:%d)/%s?%s").
	// If present (non-empty), it will be used to generate the DSN (with
	// fmt.Sprintf). You can embed placeholders for user, password, host,
	// port, database, and parameters. If left blank, the DSN() function
	// will fall back to a default per-driver build.
	DSNFormat string
}

// Connect establishes a connection to the database using the provided
// configuration.
//
//   - cfg: The configuration for the database connection.
//   - dbFactory: The factory function to create the database driver.
//   - dsn: The database connection string.
//
// Returns: the open database (DB) or an error.
func Connect(
	cfg ConnectConfig,
	dbFactory DriverFactory,
	dsn string,
) (DB, error) {
	db, err := dbFactory(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	configureConnection(db, cfg)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// DSN generates a database connection string based on the provided
// configuration.
//
//   - cfg: The configuration for the database connection.
//
// If cfg.DSNFormat is non-empty, DSN() uses that format via fmt.Sprintf.
// Else it falls back to a basic switch-based approach for known drivers.
func DSN(cfg ConnectConfig) (*string, error) {
	var dsn string

	// If a DSNFormat is explicitly set, use that first.
	if cfg.DSNFormat != "" {
		dsn = fmt.Sprintf(
			cfg.DSNFormat,
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
			cfg.Parameters,
		)
		return &dsn, nil
	}

	// Otherwise, fallback to default building per driver.
	switch cfg.Driver {
	case MySQL:
		// e.g., "user:pass@tcp(host:port)/dbname?param=value"
		connType := cfg.ConnectionType
		if connType == "" {
			connType = TCP // default
		}
		if connType == TCP {
			dsn = fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s?%s",
				cfg.User,
				cfg.Password,
				cfg.Host,
				cfg.Port,
				cfg.Database,
				cfg.Parameters,
			)
		} else if connType == Unix {
			dsn = fmt.Sprintf(
				"%s:%s@unix(%s/%s)/%s?%s",
				cfg.User,
				cfg.Password,
				cfg.SocketDirectory,
				cfg.SocketName,
				cfg.Database,
				cfg.Parameters,
			)
		} else {
			return nil, fmt.Errorf("unsupported connection type: %s", connType)
		}

	case Postgres:
		// e.g., "postgres://user:pass@host:port/dbname?param=value"
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
			cfg.Parameters,
		)

	case SQLite3:
		// e.g. "filename.db?some=params"
		// or ":memory:?some=params"
		dsn = cfg.Database
		if cfg.Parameters != "" {
			dsn += "?" + cfg.Parameters
		}

	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}

	return &dsn, nil
}

// configureConnection sets up the runtime connection limits.
func configureConnection(db DB, cfg ConnectConfig) {
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
}
