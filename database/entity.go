package database

import (
	"fmt"
)

// Insert inserts a single record into the database for the given entity.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - entity: The entity to insert (provides table name and values).
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - int64: The new record's ID if applicable (e.g., auto-increment ID).
//   - error: Any error that occurred during the insertion or error checking.
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
//   - entities: A slice of entities to insert.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - int64: The new record's ID if applicable (e.g., auto-increment ID).
//   - error: Any error that occurred during the insertion or error checking.
func InsertMany(
	preparer Preparer,
	entities []Mutator,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if len(entities) == 0 {
		return 0, nil
	}
	// Collect values from all mutators
	insertedFuncs := make([]InsertedValuesFn, len(entities))
	for i, ins := range entities {
		insertedFuncs[i] = ins.InsertedValues
	}
	tableName := entities[0].TableName()
	query, args := queryBuilder.InsertMany(tableName, insertedFuncs)
	result, err := Exec(preparer, query, args)
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
//   - errorChecker: Error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - int64: The new record's ID if applicable (e.g., auto-increment ID).
//   - error: Any error that occurred during the insertion or error checking.
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

// Get retrieves a single entity of type T from the database that matches the
// given options. The function returns an error if the entity is not found.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - options: Filter and query options for the query.
//   - factoryFn: A function that returns a new instance of T.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - T: The retrieved entity of type T.
//   - error: An error if not found or on failure.
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

// GetMany retrieves multiple entities of type T from the database that match
// the given options.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - options: Filter and query options for the query.
//   - factoryFn: A function that returns a new instance of T.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Error checker to translate SQL driver errors into
//     custom errors or skip them.
//
// Returns:
//   - T: A slice of retrieved entities of type T.
//   - error: An error if not found or on failure.
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

// Count returns the count of records for the given table matching the provided
// options.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - options: Filter and query options for the query.
//   - factoryFn: A function that returns a new instance of T.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Error checker to translate SQL driver errors into custom
//     errors or skip them.
//
// Returns:
//   - int: The count of matching records.
//   - error: An error if the query fails.
func Count[T Getter](
	preparer Preparer,
	options *CountOptions,
	factoryFn func() T,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int, error) {
	table := factoryFn().TableName()
	query, params := queryBuilder.Count(table, options)
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return 0, errorChecker.Check(err)
	}
	defer stmt.Close()
	var count int
	if err := stmt.QueryRow(params...).Scan(&count); err != nil {
		return 0, errorChecker.Check(err)
	}
	return count, nil
}

// Update applies the given field updates to all records matching the selectors.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - tableNamer: An entity or struct that provides the target table name.
//   - selectors: Conditions to match target records.
//   - queryBuilder: The SQL query builder for constructing the query.
//   - errorChecker: Error checker to translate SQL driver errors into custom
//     errors or skip them.
//
// Returns:
//   - int64: The number of updated records.
//   - error: An error if the update fails.
func Update(
	preparer Preparer,
	tableNamer TableNamer,
	selectors []Selector,
	updateFields []UpdateField,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if len(updateFields) == 0 {
		return 0, nil
	}
	query, args := queryBuilder.UpdateQuery(
		tableNamer.TableName(), updateFields, selectors,
	)
	result, err := Exec(preparer, query, args)
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
//   - errorChecker: Error checker to translate SQL driver errors into custom
//     errors or skip them.
//
// Returns:
//   - int64: The number of deleted records.
//   - error: An error if the delete fails.
func Delete(
	preparer Preparer,
	tableNamer TableNamer,
	selectors []Selector,
	opts *DeleteOptions,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	query, params := queryBuilder.Delete(
		tableNamer.TableName(), selectors, opts,
	)
	result, err := Exec(preparer, query, params)
	if err != nil {
		return 0, errorChecker.Check(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errorChecker.Check(err)
	}
	return rowsAffected, nil
}

// Exec prepares and executes an SQL query and returns the Result.
// It ensures the prepared statement is closed after execution.
//
// Parameters:
//   - preparer: The database connection or transaction to use.
//   - query: The SQL query string to execute.
//   - parameters: The query parameters.
//
// Returns:
//   - Result: The Result of execution.
//   - error: An error if the execution fails.
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
//
// Parameters:
//   - db: The database connection (must implement Exec).
//   - query: The SQL query string to execute.
//   - parameters: The query parameters.
//
// Returns:
//   - Result: The Result of execution.
//   - error: An error if the execution fails.
func ExecRaw(db DB, query string, parameters []any) (Result, error) {
	result, err := db.Exec(query, parameters...)
	if err != nil {
		return nil, err
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
//
// Returns:
//   - Rows: The rows of the query. Must be closed by the caller.
//   - Stmt: The prepared statement. Must be closed by the caller.
//   - error: An error if the execution fails.
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

// checkInsertResult checks the result of an insert or upsert operation and
// returns the new ID.
func checkInsertResult(
	result Result, err error, errorChecker ErrorChecker,
) (int64, error) {
	// Use the error checker to translate errors (e.g., duplicate key).
	if err != nil {
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
func checkUpdateResult(
	result Result, err error, errorChecker ErrorChecker,
) (int64, error) {
	// Use the error checker to translate errors (e.g., duplicate key).
	if err != nil {
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
	rows, stmt, err := Query(preparer, query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer stmt.Close()
	return rowsToEntities(rows, factoryFn)
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
