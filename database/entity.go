package database

import (
	"context"
	"database/sql"
	"fmt"
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
	Scan(dest ...any) error
	Close() error
	Err() error
}

// Row wraps *sql.Row for scanning a single result.
type Row interface {
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

// ErrorChecker translates database-specific errors into application errors.
type ErrorChecker interface {
	Check(err error) error
}

// Insert inserts a single record into the database for the given entity.
//
// Parameters:
//   - preparer: The database connection or transaction to use for preparing the statement.
//   - mutator: The entity to insert (provides table name and values).
//   - queryBuilder: The SQL query builder for constructing the INSERT statement.
//   - errorChecker: Translates SQL driver errors into *core.APIError instances.
//
// Returns: The new record's ID if applicable (e.g., auto-increment ID) and any error.
func Insert(
	preparer Preparer,
	entity Mutator,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	query, args := queryBuilder.Insert(
		entity.TableName(), entity.InsertedValues,
	)
	result, err := Exec(preparer, query, args)
	return checkInsertResult(result, err, errorChecker)
}

// InsertMany inserts multiple entities in one batch operation.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - mutators: A slice of entities to insert.
//   - queryBuilder: The SQL query builder for constructing the batch INSERT.
//   - errorChecker: Translates driver errors (e.g., duplicates) into *core.APIError.
//
// Returns: The last inserted ID from the batch (if supported) and any error.
func InsertMany(
	preparer Preparer,
	mutators []Mutator,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if len(mutators) == 0 {
		return 0, nil
	}
	// Collect values from all mutators
	insertedFuncs := make([]InsertedValuesFn, len(mutators))
	for i, ins := range mutators {
		insertedFuncs[i] = ins.InsertedValues
	}
	tableName := mutators[0].TableName()
	query, args := queryBuilder.InsertMany(tableName, insertedFuncs)
	result, err := Exec(preparer, query, args)
	return checkInsertResult(result, err, errorChecker)
}

// UpsertMany performs an "insert or update" (upsert) for multiple entities in one operation.
// This is useful for bulk inserts that should update on key conflicts (if supported by the DB).
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - mutators: The entities to upsert.
//   - updateProjections: Columns and values to update if a conflict occurs.
//   - queryBuilder: The SQL builder to construct the UPSERT statement.
//   - errorChecker: Translates driver errors to *core.APIError (e.g., unique constraint violations).
//
// Returns: The last inserted ID or 0 if no insert happened, and any error.
func UpsertMany(
	preparer Preparer,
	mutators []Mutator,
	updateProjections []Projection,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if len(mutators) == 0 {
		return 0, fmt.Errorf("must provide entities to upsert")
	}
	if len(updateProjections) == 0 {
		return 0, fmt.Errorf("must provide update projections")
	}
	if len(updateProjections[0].Alias) == 0 {
		return 0, fmt.Errorf("update projections must include an alias for the upserted table")
	}
	// Prepare batch values for upsert
	insertedFuncs := make([]InsertedValuesFn, len(mutators))
	for i, ins := range mutators {
		insertedFuncs[i] = ins.InsertedValues
	}
	query, args := queryBuilder.UpsertMany(
		mutators[0].TableName(), insertedFuncs, updateProjections,
	)
	result, err := Exec(preparer, query, args)
	return checkInsertResult(result, err, errorChecker)
}

// Get retrieves a single entity of type T from the database that matches the given options.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - options: Filter and query options (e.g., WHERE clauses, joins) for the SELECT.
//   - factoryFn: A function that returns a new instance of T (for scanning).
//   - queryBuilder: The SQL builder for constructing the SELECT statement.
//   - errorChecker: Translates not-found errors into *core.APIError (e.g., NoRows).
//
// Returns: The retrieved entity of type T, or an error if not found or on failure.
func Get[T Getter](
	preparer Preparer,
	options *GetOptions,
	factoryFn func() T,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (T, error) {
	var zero T
	query, params := queryBuilder.Get(factoryFn().TableName(), options)
	entity, err := querySingle(preparer, query, params, factoryFn)
	if err != nil {
		return zero, errorChecker.Check(err)
	}
	return entity, nil
}

// GetMany retrieves multiple entities of type T from the database that match the given options.
//
// Parameters are similar to Get, but this function returns all matching records.
//
// Returns: A slice of T entities (possibly empty if no matches), or an error on failure.
func GetMany[T Getter](
	preparer Preparer,
	options *GetOptions,
	factoryFn func() T,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) ([]T, error) {
	query, params := queryBuilder.Get(factoryFn().TableName(), options)
	entities, err := queryMultiple(preparer, query, params, factoryFn)
	if err != nil {
		return nil, errorChecker.Check(err)
	}
	return entities, nil
}

// Count returns the count of records for the given table matching the provided options.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - options: Options (filters) to apply to the COUNT query.
//   - factoryFn: A function to produce a T instance (only used to get the table name).
//   - queryBuilder: The SQL builder for constructing the COUNT query.
//
// Returns: The number of records matching the criteria, or an error.
func Count[T Getter](
	preparer Preparer,
	options *CountOptions,
	factoryFn func() T,
	queryBuilder QueryBuilder,
) (int, error) {
	table := factoryFn().TableName()
	query, params := queryBuilder.Count(table, options)
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var count int
	if err := stmt.QueryRow(params...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Update applies the given field updates to all records matching the selectors.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - updater: An entity or struct that provides the target table name.
//   - selectors: Conditions to match target records (e.g., WHERE clauses).
//   - updateFields: List of fields (column names and new values) to update.
//   - queryBuilder: The SQL builder for constructing the UPDATE statement.
//   - errorChecker: Translates driver errors (if any) via Check.
//
// Returns: The number of rows affected, and any error.
func Update(
	preparer Preparer,
	updater TableNamer,
	selectors []Selector,
	updateFields []UpdateField,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if len(updateFields) == 0 {
		return 0, nil
	}
	query, args := queryBuilder.UpdateQuery(
		updater.TableName(), updateFields, selectors,
	)
	result, err := Exec(preparer, query, args)
	return checkUpdateResult(result, err, errorChecker)
}

// Delete removes records from the database table matching the given selectors.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - deleter: Provides the table name (usually an empty instance or struct of the entity).
//   - selectors: Conditions to identify records to delete (e.g., primary key equals X).
//   - queryBuilder: The SQL builder for constructing the DELETE statement.
//   - opts: Additional options for the delete (e.g., LIMIT).
//
// Returns: The number of rows deleted, or an error.
func Delete(
	preparer Preparer,
	deleter TableNamer,
	selectors []Selector,
	opts *DeleteOptions,
	queryBuilder QueryBuilder,
) (int64, error) {
	query, params := queryBuilder.Delete(deleter.TableName(), selectors, opts)
	result, err := Exec(preparer, query, params)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// Exec prepares and executes a non-query SQL statement (INSERT, UPDATE, DELETE, etc.).
// It ensures the prepared statement is closed after execution.
//
// Parameters:
//   - preparer: The database or transaction to use for preparing the statement.
//   - query: The SQL query string to execute.
//   - parameters: The query parameters.
//
// Returns: The Result of the execution (for retrieving last insert ID, rows affected) and an error if execution fails.
func Exec(preparer Preparer, query string, parameters []any) (Result, error) {
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(parameters...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ExecRaw executes a query directly on the DB without explicit preparation.
// This can be used for one-off statements where parameterization is handled by the driver.
//
// Parameters:
//   - db: The database connection (must implement Exec).
//   - query: The SQL query string to execute.
//   - parameters: The query parameters.
//
// Returns: The Result of execution and any error.
func ExecRaw(db DB, query string, parameters []any) (Result, error) {
	result, err := db.Exec(query, parameters...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Query prepares and executes a query that returns rows. It returns both the Rows and the Stmt.
// The caller is responsible for closing both the Rows and the Stmt when done.
//
// Parameters:
//   - preparer: The database or transaction to use for preparing the statement.
//   - query: The SQL SELECT query string.
//   - parameters: The query parameters.
//
// Returns: The resulting Rows and the prepared Stmt (both must be closed by the caller), or an error.
func Query(
	preparer Preparer, query string, parameters []any,
) (Rows, Stmt, error) {
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return nil, nil, err
	}
	rows, err := stmt.Query(parameters...)
	if err != nil {
		stmt.Close()
		return nil, nil, err
	}
	return rows, stmt, nil
}

// checkInsertResult checks the result of an insert or upsert operation and returns the new ID.
func checkInsertResult(
	result Result, err error, errorChecker ErrorChecker,
) (int64, error) {
	if err != nil {
		// Use the error checker to translate errors (e.g., duplicate key) if possible.
		return 0, errorChecker.Check(err)
	}
	if result == nil {
		return 0, nil // No result (no ID available).
	}
	id, err := result.LastInsertId()
	if err != nil {
		// If LastInsertId isn't supported, it's not a fatal error; return 0 ID.
		return 0, nil
	}
	return id, nil
}

// checkUpdateResult checks the result of an update and returns rows affected.
func checkUpdateResult(
	result Result, err error, errorChecker ErrorChecker,
) (int64, error) {
	if err != nil {
		return 0, errorChecker.Check(err)
	}
	if result == nil {
		return 0, nil
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// queryMultiple is an internal helper to query and scan multiple entities of type T.
func queryMultiple[T Getter](
	preparer Preparer, query string, params []any, factoryFn func() T,
) ([]T, error) {
	rows, stmt, err := Query(preparer, query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer stmt.Close()
	return rowsToEntities(rows, factoryFn)
}

// querySingle is an internal helper to query and scan a single entity of type T.
func querySingle[T Getter](
	preparer Preparer, query string, params []any, factoryFn func() T,
) (T, error) {
	var zero T
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return zero, err
	}
	defer stmt.Close()
	// Use QueryRow since we expect at most one result.
	return rowToEntity(stmt.QueryRow(params...), factoryFn)
}

// rowToEntity scans a single Row into a new entity of type T.
func rowToEntity[T Getter](row Row, factoryFn func() T) (T, error) {
	var zero T
	entity := factoryFn()
	if err := entity.ScanRow(row); err != nil {
		return zero, err
	}
	if err := row.Err(); err != nil {
		return zero, err
	}
	return entity, nil
}

// rowsToEntities scans all rows into a slice of entities of type T.
func rowsToEntities[T Getter](rows Rows, factoryFn func() T) ([]T, error) {
	results := []T{}
	for rows.Next() {
		entity := factoryFn()
		if err := entity.ScanRow(rows); err != nil {
			return nil, err
		}
		results = append(results, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
