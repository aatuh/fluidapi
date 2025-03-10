package database

import (
	"fmt"
	"time"
)

// Supported drivers.
const (
	MySQL      = "mysql"      // MySQL driver name
	PostgreSQL = "postgresql" // PostgreSQL driver name
	SQLite3    = "sqlite3"    // SQLite3 driver name
)

// Supported connection types.
const (
	TCP  = "tcp"  // TCP connection type
	Unix = "unix" // Unix socket connection type
)

// ConnectConfig holds the configuration for the database connection.
type ConnectConfig struct {
	User            string        // Database user
	Password        string        // Database password
	Host            string        // Database host
	Port            int           // Database port
	Database        string        // Database name (e.g. "users")
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

// DriverFactory is a function that creates a database driver.
type DriverFactory func(driver string, dsn string) (DB, error)

// Connect establishes a connection to the database using the provided
// configuration. It will automatically configure the connection based on the
// provided configuration and then attempt to ping the database.
//
// Parameters:
//   - cfg: The configuration for the database connection.
//   - dbFactory: The factory function to create the database driver.
//   - dsn: The database connection string.
//
// Returns:
//   - DB: The database connection.
//   - error: An error if the connection fails.
func Connect(
	cfg ConnectConfig, dbFactory DriverFactory, dsn string,
) (DB, error) {
	db, err := dbFactory(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("Connect: failed to open database: %w", err)
	}
	configureConnection(db, cfg)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Connect: failed to ping database: %w", err)
	}
	return db, nil
}

// DSN generates a database connection string based on the provided
// configuration. If cfg.DSNFormat is non-empty, DSN() uses that format via
// fmt.Sprintf. Else it falls back to determining DSN for known drivers.
//
// Parameters:
//   - cfg: The configuration for the database connection.
//
// Returns:
//   - *string: A pointer to the generated DSN string.
//   - error: An error if the DSN generation fails.
func DSN(cfg ConnectConfig) (*string, error) {
	var dsn string
	// If DSNFormat is set, use it.
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
	// Otherwise, build the DSN based on the driver.
	var err error
	switch cfg.Driver {
	case MySQL:
		dsn, err = dsnForMySQL(cfg)
	case PostgreSQL:
		dsn, err = dsnForPostgreSQL(cfg)
	case SQLite3:
		dsn, err = dsnForSQLite3(cfg)
	default:
		return nil, fmt.Errorf("DSN: unsupported driver: %s", cfg.Driver)
	}
	if err != nil {
		return nil, fmt.Errorf("DSN: failed to generate DSN: %w", err)
	}
	return &dsn, nil
}

// dsnForMySQL builds a DSN for MySQL.
func dsnForMySQL(cfg ConnectConfig) (string, error) {
	connType := cfg.ConnectionType
	if connType == "" {
		connType = TCP
	}
	switch connType {
	case TCP:
		return dsnForMySQLTCP(cfg)
	case Unix:
		return dsnForMySQLUnix(cfg)
	default:
		return "", fmt.Errorf(
			"dsnForMySQL: unsupported connection type: %s", connType,
		)
	}
}

// dsnForMySQLTCP builds a DSN for MySQL over TCP socket.
func dsnForMySQLTCP(cfg ConnectConfig) (dsn string, err error) {
	if cfg.Host == "" {
		return "", fmt.Errorf("dsnForMySQL: missing required MySQL host field")
	}
	if cfg.Port == 0 {
		return "", fmt.Errorf("dsnForMySQL: missing required MySQL port field")
	}
	if cfg.Database == "" {
		return "", fmt.Errorf("dsnForMySQL: missing required MySQL database field")
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Parameters)
	return dsn, nil
}

// dsnForMySQLUnix builds a DSN for MySQL over Unix socket.
func dsnForMySQLUnix(cfg ConnectConfig) (string, error) {
	if cfg.SocketDirectory == "" {
		return "", fmt.Errorf("dsnForMySQL: missing required MySQL socket directory field")
	}
	if cfg.SocketName == "" {
		return "", fmt.Errorf("dsnForMySQL: missing required MySQL socket name field")
	}
	if cfg.Database == "" {
		return "", fmt.Errorf("dsnForMySQL: missing required MySQL database field")
	}

	return fmt.Sprintf(
		"%s:%s@unix(%s/%s)/%s?%s",
		cfg.User,
		cfg.Password,
		cfg.SocketDirectory,
		cfg.SocketName,
		cfg.Database,
		cfg.Parameters,
	), nil
}

// dsnForPostgres builds a DSN for PostgreSQL.
func dsnForPostgreSQL(cfg ConnectConfig) (string, error) {
	if cfg.Host == "" {
		return "", fmt.Errorf(
			"dsnForPostgreSQL: missing required PostgreSQL host field",
		)
	}
	if cfg.Port == 0 {
		return "", fmt.Errorf(
			"dsnForPostgreSQL: missing required PostgreSQL port field",
		)
	}
	if cfg.Database == "" {
		return "", fmt.Errorf(
			"dsnForPostgreSQL: missing required PostgreSQL database field",
		)
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Parameters,
	), nil
}

// dsnForSQLite3 builds a DSN for SQLite3.
func dsnForSQLite3(cfg ConnectConfig) (string, error) {
	if cfg.Database == "" {
		return "", fmt.Errorf(
			"dsnForSQLite3: database name is required for SQLite3",
		)
	}
	dsn := cfg.Database
	if cfg.Parameters != "" {
		dsn += "?" + cfg.Parameters
	}
	return dsn, nil
}

// configureConnection sets up the runtime connection limits.
func configureConnection(db DB, cfg ConnectConfig) {
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
}
