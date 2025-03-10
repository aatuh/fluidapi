package database

import (
	"context"
	"database/sql"
	"time"
)

// Preparer is an interface for preparing SQL statements.
type Preparer interface {
	Prepare(query string) (Stmt, error)
}

// DB is an interface for core database operations and connection management.
type DB interface {
	Preparer
	Ping() error
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Exec(query string, args ...any) (Result, error)
	Query(query string, args ...any) (Rows, error)
	Close() error
}

// Tx is an interface for transaction operations.
type Tx interface {
	Preparer
	Commit() error
	Rollback() error
	Exec(query string, args ...any) (Result, error)
}

// Stmt wraps *sql.Stmt methods for executing prepared statements.
type Stmt interface {
	Close() error
	QueryRow(args ...any) Row
	Exec(args ...any) (Result, error)
	Query(args ...any) (Rows, error)
}

// Rows wraps *sql.Rows for scanning multiple results.
type Rows interface {
	Next() bool
	// Scan scans the rows into dest values.
	Scan(dest ...any) error
	Close() error
	Err() error
}

// Row wraps *sql.Row for scanning a single result.
type Row interface {
	// Scan scans the row into dest values.
	Scan(dest ...any) error
	Err() error
}

// Result wraps *sql.Result for retrieving metadata of write operations.
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// TableNamer provides the table name for an entity.
type TableNamer interface {
	TableName() string
}

// Mutator provides values for insert and update operations.
type Mutator interface {
	TableNamer
	// InsertedValues returns the column names and values for insertion.
	InsertedValues() ([]string, []any)
}

// Getter can scan a database row into itself.
type Getter interface {
	TableNamer
	// ScanRow should populate the entity from the given Row.
	ScanRow(row Row) error
}

// CRUDEntity is a helper constraint for entities that can be inserted,
// retrieved, and updated.
type CRUDEntity interface {
	Mutator
	Getter
}

// ErrorChecker translates database-specific errors into application errors.
type ErrorChecker interface {
	Check(err error) error
}
