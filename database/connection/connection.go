package connection

import (
	"fmt"
	"time"

	"github.com/pakkasys/fluidapi/database/util"
)

const (
	TCP  = "tcp"  // TCP connection type
	Unix = "unix" // Unix socket connection type

	MySQL    = "mysql"    // MySQL driver name
	Postgres = "postgres" // PostgreSQL driver name
	SQLite3  = "sqlite3"  // SQLite3 driver name
)

// DriverFactory is a function that creates a database driver.
type DriverFactory func(driver string, dsn string) (util.DB, error)

// Config holds the configuration for the database connection.
type Config struct {
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
	DSNFormat       string        // Custom DSN format
}

// Connect establishes a connection to the database using the provided
// configuration.
//
//   - cfg: The configuration for the database connection.
//   - dbFactory: The factory function to create the database driver.
//   - dsn: The database connection string.
func Connect(
	cfg Config,
	dbFactory DriverFactory,
	dsn string,
) (util.DB, error) {
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

// GetDSN generates a database connection string based on the provided
// configuration.
//
//   - cfg: The configuration for the database connection.
func GetDSN(cfg Config) (*string, error) {
	var dsn string
	switch cfg.Driver {
	case MySQL:
		// MySQL DSN is usually something like:
		// "user:pass@tcp(host:port)/dbname?param=value"
		dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?%s",
			cfg.User,
			cfg.Password,
			cfg.ConnectionType, // "tcp" or "unix"
			cfg.Host,
			cfg.Port,
			cfg.Database,
			cfg.Parameters,
		)
	case Postgres:
		// Postgres DSN is usually something like:
		// "postgres://user:pass@host:port/dbname?param=value"
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
			cfg.Parameters,
		)
	case SQLite3:
		// SQLite DSN is usually just the file path + params
		dsn = fmt.Sprintf("%s?%s", cfg.Database, cfg.Parameters)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}
	return &dsn, nil
}

func configureConnection(db util.DB, cfg Config) {
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
}
