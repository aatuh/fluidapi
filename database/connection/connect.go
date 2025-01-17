package connection

import (
	"fmt"
	"time"

	"github.com/pakkasys/fluidapi/database/util"
)

const (
	TCP  = "tcp"  // TCP connection type
	Unix = "unix" // Unix socket connection type
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
func Connect(cfg *Config, dbFactory DriverFactory) (util.DB, error) {
	dsn, err := getDSN(cfg)
	if err != nil {
		return nil, err
	}

	db, err := dbFactory(cfg.Driver, *dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	configureConnection(db, cfg)

	return db, nil
}

func getDSN(cfg *Config) (*string, error) {
	switch cfg.ConnectionType {
	case TCP:
		dsn := fmt.Sprintf(
			cfg.DSNFormat,
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
			cfg.Parameters,
		)
		return &dsn, nil
	case Unix:
		dsn := fmt.Sprintf(
			cfg.DSNFormat,
			cfg.User,
			cfg.Password,
			cfg.SocketDirectory,
			cfg.SocketName,
			cfg.Database,
			cfg.Parameters,
		)
		return &dsn, nil
	default:
		return nil, fmt.Errorf(
			"unsupported connection type: %s",
			cfg.ConnectionType,
		)
	}
}

// TODO: Opinionated, move elsewhere
func configureConnection(db util.DB, cfg *Config) {
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
}
