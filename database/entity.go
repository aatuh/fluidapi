package database

import (
	"fmt"
)

// ReadDBOps provides methods to perform database read operations.
type ReadDBOps[Entity Getter] struct{}

// NewReadDBOps creates a new ReadDBOps instance.
func NewReadDBOps[Entity Getter]() *ReadDBOps[Entity] {
	return &ReadDBOps[Entity]{}
}

// Get retrieves a single entity of type T from the database that matches the
// given options. The function returns an error if the entity is not found.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - options: Filter and query options for the query.
//   - factoryFn: A function that returns a new instance of T.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - T: The retrieved entity of type T.
//   - error: An error if not found or on failure.
func (d *ReadDBOps[Entity]) Get(
	preparer Preparer,
	options *GetOptions,
	factoryFn func() Entity,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (Entity, error) {
	var zero Entity
	if preparer == nil {
		return zero, fmt.Errorf("Get: preparer is nil")
	}
	if options == nil {
		return zero, fmt.Errorf("Get: options is nil")
	}
	if queryBuilder == nil {
		return zero, fmt.Errorf("Get: queryBuilder is nil")
	}

	query, params := queryBuilder.Get(factoryFn().TableName(), options)
	entity, err := querySingle(preparer, query, params, factoryFn)
	if err != nil {
		if errorChecker == nil {
			return zero, err
		}
		return zero, errorChecker.Check(err)
	}
	return entity, nil
}

// GetMany retrieves multiple entities of type T from the database that match
// the given options.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - options: Filter and query options for the query.
//   - factoryFn: A function that returns a new instance of T.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - T: A slice of retrieved entities of type T.
//   - error: An error if not found or on failure.
func (d *ReadDBOps[Entity]) GetMany(
	preparer Preparer,
	options *GetOptions,
	factoryFn func() Entity,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) ([]Entity, error) {
	if preparer == nil {
		return nil, fmt.Errorf("GetMany: preparer is nil")
	}
	if options == nil {
		return nil, fmt.Errorf("GetMany: options is nil")
	}
	if factoryFn == nil {
		return nil, fmt.Errorf("GetMany: factoryFn is nil")
	}
	if queryBuilder == nil {
		return nil, fmt.Errorf("GetMany: queryBuilder is nil")
	}

	query, params := queryBuilder.Get(factoryFn().TableName(), options)
	entities, err := queryMultiple(preparer, query, params, factoryFn)
	if err != nil {
		if errorChecker == nil {
			return nil, err
		}
		return nil, errorChecker.Check(err)
	}
	return entities, nil
}

// Count returns the count of records for the given table matching the provided
// options.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - options: Filter and query options for the query.
//   - factoryFn: A function that returns a new instance of T.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - int: The count of matching records.
//   - error: An error if the query fails.
func (d *ReadDBOps[Entity]) Count(
	preparer Preparer,
	options *CountOptions,
	factoryFn func() Entity,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int, error) {
	if preparer == nil {
		return 0, fmt.Errorf("Count: preparer is nil")
	}
	if options == nil {
		return 0, fmt.Errorf("Count: options is nil")
	}
	if factoryFn == nil {
		return 0, fmt.Errorf("Count: factoryFn is nil")
	}
	if queryBuilder == nil {
		return 0, fmt.Errorf("Count: queryBuilder is nil")
	}

	table := factoryFn().TableName()
	query, params := queryBuilder.Count(table, options)
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return 0, errorChecker.Check(err)
	}
	defer stmt.Close()
	var count int
	if err := stmt.QueryRow(params...).Scan(&count); err != nil {
		if errorChecker == nil {
			return 0, err
		}
		return 0, errorChecker.Check(err)
	}
	return count, nil
}

// ReadDBOps provides methods to perform database write operations.
type MutateDBOps[Entity Mutator] struct{}

// NewMutateDBOps creates a new MutateDBOps instance.
func NewMutateDBOps[Entity Mutator]() *MutateDBOps[Entity] {
	return &MutateDBOps[Entity]{}
}

// Insert inserts a single record into the database for the given entity.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - entity: The entity to insert (provides table name and values).
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - int64: The new record's ID if applicable (e.g., auto-increment ID).
//   - error: Any error that occurred during the insertion or error checking.
func (d *MutateDBOps[Entity]) Insert(
	preparer Preparer,
	entity Mutator,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if preparer == nil {
		return 0, fmt.Errorf("Insert: preparer is nil")
	}
	if entity == nil {
		return 0, fmt.Errorf("Insert: entity is nil")
	}
	if queryBuilder == nil {
		return 0, fmt.Errorf("Insert: queryBuilder is nil")
	}

	query, args := queryBuilder.Insert(
		entity.TableName(), entity.InsertedValues,
	)
	result, err := doExec(preparer, query, args)
	return checkInsertResult(result, err, errorChecker)
}

// InsertMany inserts multiple entities in one batch operation.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - entities: A slice of entities to insert.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - int64: The new record's ID if applicable (e.g., auto-increment ID).
//   - error: Any error that occurred during the insertion or error checking.
func (d *MutateDBOps[Entity]) InsertMany(
	preparer Preparer,
	entities []Mutator,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if preparer == nil {
		return 0, fmt.Errorf("InsertMany: preparer is nil")
	}
	if len(entities) == 0 {
		return 0, nil
	}
	if queryBuilder == nil {
		return 0, fmt.Errorf("InsertMany: queryBuilder is nil")
	}

	// Collect values from all mutators
	insertedFuncs := make([]InsertedValuesFn, len(entities))
	for i, ins := range entities {
		insertedFuncs[i] = ins.InsertedValues
	}
	tableName := entities[0].TableName()
	query, args := queryBuilder.InsertMany(tableName, insertedFuncs)
	result, err := doExec(preparer, query, args)
	return checkInsertResult(result, err, errorChecker)
}

// UpsertMany performs an "insert or update" (upsert) for multiple entities in
// one operation. This is useful for bulk inserts that should update on key
// conflicts (if supported by the DB).
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - entities: A slice of entities to upsert.
//   - updateProjections: Columns and values to update if a conflict occurs.
//   - entity: The entity to insert (provides table name and values).
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - int64: The new record's ID if applicable (e.g., auto-increment ID).
//   - error: Any error that occurred during the insertion or error checking.
func (d *MutateDBOps[Entity]) UpsertMany(
	preparer Preparer,
	mutators []Mutator,
	updateProjections []Projection,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if preparer == nil {
		return 0, fmt.Errorf("UpsertMany: preparer is nil")
	}
	if len(mutators) == 0 {
		return 0, fmt.Errorf("UpsertMany: must provide entities to upsert")
	}
	if len(updateProjections) == 0 {
		return 0, fmt.Errorf("UpsertMany: must provide update projections")
	}
	if len(updateProjections[0].Alias) == 0 {
		return 0, fmt.Errorf("UpsertMany: update projections must include an alias for the upserted table")
	}

	// Prepare batch values for upsert
	insertedFuncs := make([]InsertedValuesFn, len(mutators))
	for i, ins := range mutators {
		insertedFuncs[i] = ins.InsertedValues
	}
	query, args := queryBuilder.UpsertMany(
		mutators[0].TableName(), insertedFuncs, updateProjections,
	)
	result, err := doExec(preparer, query, args)
	return checkInsertResult(result, err, errorChecker)
}

// Update applies the given field updates to all records matching the selectors.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - tableNamer: An entity or struct that provides the target table name.
//   - selectors: Conditions to match target records.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - int64: The number of updated records.
//   - error: An error if the update fails.
func (*MutateDBOps[Entity]) Update(
	preparer Preparer,
	tableNamer TableNamer,
	selectors []Selector,
	updates []Update,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if preparer == nil {
		return 0, fmt.Errorf("Update: preparer is nil")
	}
	if tableNamer == nil {
		return 0, fmt.Errorf("Update: tableNamer is nil")
	}
	if len(updates) == 0 {
		return 0, nil
	}
	if queryBuilder == nil {
		return 0, fmt.Errorf("Update: queryBuilder is nil")
	}

	query, args := queryBuilder.UpdateQuery(
		tableNamer.TableName(), updates, selectors,
	)
	result, err := doExec(preparer, query, args)
	return checkUpdateResult(result, err, errorChecker)
}

// Delete removes records from the database table matching the given selectors.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - tableNamer: An entity or struct that provides the target table name.
//   - selectors: Conditions to match target records.
//   - opts: Options for the delete operation.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Optional error checker to translate SQL driver errors into custom
//     errors or skip them.
//
// Returns:
//   - int64: The number of deleted records.
//   - error: An error if the delete fails.
func (*MutateDBOps[Entity]) Delete(
	preparer Preparer,
	entity Mutator,
	selectors []Selector,
	opts *DeleteOptions,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if preparer == nil {
		return 0, fmt.Errorf("Delete: preparer is nil")
	}
	if entity == nil {
		return 0, fmt.Errorf("Delete: tableNamer is nil")
	}
	if opts == nil {
		return 0, fmt.Errorf("Delete: opts is nil")
	}
	if queryBuilder == nil {
		return 0, fmt.Errorf("Delete: queryBuilder is nil")
	}

	query, params := queryBuilder.Delete(
		entity.TableName(), selectors, opts,
	)
	result, err := doExec(preparer, query, params)
	if err != nil {
		return 0, errorChecker.Check(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if errorChecker == nil {
			return 0, err
		}
		return 0, errorChecker.Check(err)
	}
	return rowsAffected, nil
}

// DBOps provides methods to perform database exec and query operations.
type DBOps struct{}

// NewDBOps creates a new DBOps instance.
func NewDBOps() *DBOps {
	return &DBOps{}
}

// Exec prepares and executes an SQL query and returns the Result.
// It ensures the prepared statement is closed after execution.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - query: The SQL query string to execute.
//   - parameters: The query parameters.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - Result: The Result of execution.
//   - error: An error if the execution fails.
func (*DBOps) Exec(
	preparer Preparer,
	query string,
	parameters []any,
	errorChecker ErrorChecker,
) (Result, error) {
	if preparer == nil {
		return nil, fmt.Errorf("Exec: preparer is nil")
	}

	result, err := doExec(preparer, query, parameters)
	if err != nil {
		if errorChecker == nil {
			return nil, err
		}
		return nil, errorChecker.Check(err)
	}
	return result, nil
}

// ExecRaw executes a query directly on the DB without explicit preparation.
//
// Parameters:
//   - db: The database connection (must implement Exec).
//   - query: The SQL query string to execute.
//   - parameters: The query parameters.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - Result: The Result of execution.
//   - error: An error if the execution fails.
func (*DBOps) ExecRaw(
	db DB, query string, parameters []any, errorChecker ErrorChecker,
) (Result, error) {
	if db == nil {
		return nil, fmt.Errorf("ExecRaw: db is nil")
	}

	result, err := doExecRaw(db, query, parameters)
	if err != nil {
		if errorChecker == nil {
			return nil, err
		}
		return nil, errorChecker.Check(err)
	}
	return result, nil
}

// Query prepares and executes a query that returns rows. It returns both the
// Rows and the Stmt. The caller is responsible for closing both the Rows and
// the Stmt when done.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - query: The SQL query string to execute.
//   - parameters: The query parameters.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - Rows: The rows of the query. Must be closed by the caller.
//   - Stmt: The prepared statement. Must be closed by the caller.
//   - error: An error if the execution fails.
func (*DBOps) Query(
	preparer Preparer,
	query string,
	parameters []any,
	errorChecker ErrorChecker,
) (Rows, Stmt, error) {
	if preparer == nil {
		return nil, nil, fmt.Errorf("Query: preparer is nil")
	}

	rows, stmt, err := doQuery(preparer, query, parameters)
	if err != nil {
		if errorChecker == nil {
			return nil, nil, err
		}
		return nil, nil, errorChecker.Check(err)
	}
	return rows, stmt, nil
}

// QueryRaw executes a query directly on the DB without explicit preparation.
//
// Parameters:
//   - db: The database connection (must implement Query).
//   - query: The SQL query string to execute.
//   - parameters: The query parameters.
//   - errorChecker: Optional error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - Rows: The rows of the query. Must be closed by the caller.
//   - error: An error if the execution fails.
func (*DBOps) QueryRaw(
	db DB, query string, parameters []any, errorChecker ErrorChecker,
) (Rows, error) {
	if db == nil {
		return nil, fmt.Errorf("QueryRaw: db is nil")
	}

	rows, err := doQueryRaw(db, query, parameters)
	if err != nil {
		if errorChecker == nil {
			return nil, err
		}
		return nil, errorChecker.Check(err)
	}
	return rows, nil
}

// RowToEntity scans a single Row into a new entity of type T.
//
// Parameters:
//   - row: The Row to scan.
//   - factoryFn: A function that returns a new instance of T.
//
// Returns:
//   - T: The scanned entity of type T.
//   - error: An error if the scan fails.
func RowToEntity[T Getter](row Row, factoryFn func() T) (T, error) {
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

// RowsToEntities scans all rows into a slice of entities of type T.
//
// Parameters:
//   - rows: The Rows to scan.
//   - factoryFn: A function that returns a new instance of T.
//
// Returns:
//   - []T: A slice of scanned entities of type T.
//   - error: An error if the scan fails.
func RowsToEntities[T Getter](rows Rows, factoryFn func() T) ([]T, error) {
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

// checkInsertResult checks the result of an insert or upsert operation and
// returns the new ID. Error checker is optional.
func checkInsertResult(
	result Result, err error, errorChecker ErrorChecker,
) (int64, error) {
	// Use the error checker to translate errors (e.g., duplicate key).
	if err != nil && errorChecker != nil {
		return 0, errorChecker.Check(err)
	}
	if result == nil {
		return 0, nil // No result (no ID available).
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errorChecker.Check(err)
	}
	return id, nil
}

// checkUpdateResult checks the result of an update and returns rows affected.
// Error checker is optional.
func checkUpdateResult(
	result Result, err error, errorChecker ErrorChecker,
) (int64, error) {
	// Use the error checker to translate errors (e.g., duplicate key).
	if err != nil && errorChecker != nil {
		return 0, errorChecker.Check(err)
	}
	if result == nil {
		return 0, nil
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, errorChecker.Check(err)
	}
	return count, nil
}

// queryMultiple queries and scans multiple entities of type T.
func queryMultiple[T Getter](
	preparer Preparer, query string, params []any, factoryFn func() T,
) ([]T, error) {
	rows, stmt, err := doQuery(preparer, query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer stmt.Close()
	return RowsToEntities(rows, factoryFn)
}

// querySingle queries and scans a single entity of type T.
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
	return RowToEntity(stmt.QueryRow(params...), factoryFn)
}

// doExec is a helper to execute an SQL query without error checking.
func doExec(preparer Preparer, query string, parameters []any) (Result, error) {
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

// doExecRaw is a helper to execute an SQL query without error checking.
func doExecRaw(db DB, query string, parameters []any) (Result, error) {
	result, err := db.Exec(query, parameters...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// doQuery is a helper to execute a query without error checking.
func doQuery(
	preparer Preparer, query string, parameters []any,
) (Rows, Stmt, error) {
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return nil, nil, err
	}
	rows, err := stmt.Query(parameters...)
	if err != nil {
		if closeErr := stmt.Close(); closeErr != nil {
			return nil, nil, fmt.Errorf(
				"query error: %w; additionally, stmt.Close error: %v",
				err,
				closeErr,
			)
		}
		return nil, nil, err
	}
	return rows, stmt, nil
}

// doQueryRaw is a helper to execute a query without error checking.
func doQueryRaw(db DB, query string, parameters []any) (Rows, error) {
	rows, err := db.Query(query, parameters...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
