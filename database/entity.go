package database

import (
	"database/sql"
	"fmt"
)

// ErrorChecker is an interface for checking database errors for a specific
// driver.
type ErrorChecker interface {
	Check(err error) error
}

// Inserter is an interface for entities that can be inserted into the
type Inserter interface {
	TableNamer
	InsertedValues() ([]string, []any)
}

// Getter is our interface for any entity that has a table name and can scan a row.
type Getter interface {
	TableNamer
	ScanRow(row Row) error
}

// Deleter is an interface for entities that can be updated in the
type TableNamer interface {
	TableName() string
}

// Insert creates an entity in the
//
//   - preparer: The database connection.
//   - inserter: The entity to insert.
//   - queryBuilder: The query builder.
//   - errorChecker: The function used to check the error.
func Insert(
	preparer Preparer,
	inserter Inserter,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	query, values := queryBuilder.Insert(
		inserter.TableName(),
		inserter.InsertedValues,
	)
	result, err := Exec(preparer, query, values)
	return checkInsertResult(result, err, errorChecker)
}

// InsertMany creates many entities in the
//
//   - preparer: The database connection.
//   - inserters: The entities to insert.
//   - queryBuilder: The query builder.
//   - errorChecker: The function used to check the error.
func InsertMany(
	preparer Preparer,
	inserters []Inserter,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if len(inserters) == 0 {
		return 0, nil
	}
	insertedValues := make([]InsertedValues, len(inserters))
	for i := range inserters {
		insertedValues[i] = inserters[i].InsertedValues
	}
	tableName := inserters[0].TableName()
	query, values := queryBuilder.InsertMany(tableName, insertedValues)
	result, err := Exec(preparer, query, values)
	return checkInsertResult(result, err, errorChecker)
}

// UpsertMany upserts multiple entities.
//
//   - db: The database connection.
//   - inserters: The entities to upsert.
//   - updateProjections: The projections of the entities to update.
//   - queryBuilder: The query builder.
//   - errorChecker: The function used to check the error.
func UpsertMany(
	preparer Preparer,
	inserters []Inserter,
	updateProjections []Projection,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (int64, error) {
	if len(inserters) == 0 {
		return 0, fmt.Errorf("must provide entities to upsert")
	}
	if len(updateProjections) == 0 {
		return 0, fmt.Errorf("must provide update projections")
	}
	if len(updateProjections[0].Alias) == 0 {
		return 0, fmt.Errorf("must provide update projections alias")
	}
	insertedValues := make([]InsertedValues, len(inserters))
	for i := range inserters {
		insertedValues[i] = inserters[i].InsertedValues
	}
	query, values := queryBuilder.UpsertMany(
		inserters[0].TableName(),
		insertedValues,
		updateProjections,
	)
	result, err := Exec(preparer, query, values)
	return checkInsertResult(result, err, errorChecker)
}

// Get returns a single entity of type T.
//
//   - preparer: The database connection.
//   - dbOptions: The options for the query.
//   - factoryFn: A function that returns a new instance of T.
//   - queryBuilder: The query builder.
//   - errorChecker: The function used to check the error.
func Get[T Getter](
	preparer Preparer,
	dbOptions *GetOptions,
	factoryFn func() T,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) (T, error) {
	queryStr, whereValues := queryBuilder.Get(
		factoryFn().TableName(),
		dbOptions,
	)
	entity, err := querySingle(preparer, queryStr, whereValues, factoryFn)
	return entity, errorChecker.Check(err)
}

// GetMany returns multiple entities of type T.
//
//   - preparer: The database connection.
//   - dbOptions: The options for the query.
//   - factoryFn: A function that returns a new instance of T.
//   - queryBuilder: The query builder.
//   - errorChecker: The function used to check the error.
func GetMany[T Getter](
	preparer Preparer,
	dbOptions *GetOptions,
	factoryFn func() T,
	queryBuilder QueryBuilder,
	errorChecker ErrorChecker,
) ([]T, error) {
	queryStr, whereValues := queryBuilder.Get(
		factoryFn().TableName(),
		dbOptions,
	)
	entities, err := queryMultiple(
		preparer,
		queryStr,
		whereValues,
		factoryFn,
	)
	return entities, errorChecker.Check(err)
}

// Count counts the number of entities in the
//
//   - preparer: The database connection.
//   - tableName: The name of the database table.
//   - queryBuilder: The query builder.
//   - dbOptions: The options for the query.
func Count[T Getter](
	preparer Preparer,
	dbOptions *CountOptions,
	factoryFn func() T,
	queryBuilder QueryBuilder,
) (int, error) {
	obj := *new(T)
	query, whereValues := queryBuilder.Count(obj.TableName(), dbOptions)

	statement, err := preparer.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var count int
	if err := statement.QueryRow(whereValues...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// Update updates entities in the
//
//   - db: The database connection to use.
//   - updater: The updater to use.
//   - selectors: The selectors of the entities to update.
//   - updateFields The fields and values to update.
//   - queryBuilder: The query builder.
//   - updates: The updates to apply to the entities.
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

	query, values := queryBuilder.UpdateQuery(
		updater.TableName(),
		updateFields,
		selectors,
	)
	result, err := Exec(preparer, query, values)
	return checkUpdateResult(result, err, errorChecker)
}

// Delete deletes entities from the
//
//   - preparer: The database connection.
//   - deleter: The deleter to use.
//   - selectors: The selectors for the entities to delete.
//   - queryBuilder: The query builder.
//   - opts: The options for the query.
func Delete(
	preparer Preparer,
	deleter TableNamer,
	selectors []Selector,
	opts *DeleteOptions,
	queryBuilder QueryBuilder,
) (int64, error) {
	query, whereValues := queryBuilder.Delete(
		deleter.TableName(),
		selectors,
		opts,
	)

	result, err := Exec(preparer, query, whereValues)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Exec runs a query and returns the result.
//
//   - preparer: The preparer used to prepare the query.
//   - query: The query string.
//   - parameters: The parameters for the query.
func Exec(
	preparer Preparer,
	query string,
	parameters []any,
) (Result, error) {
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

// ExecRaw runs a non-prepared statement query and returns the result.
//
//   - db: The database connection.
//   - query: The query string.
//   - parameters: The parameters for the query.
func ExecRaw(
	db DB,
	query string,
	parameters []any,
) (Result, error) {
	result, err := db.Exec(query, parameters...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Query runs a query and returns the rows. The returned rows and stmt objects
// must be managed manually by the caller after successful query execution.
//
//   - preparer: The preparer used to prepare the query.
//   - query: The query string.
//   - parameters: The parameters for the query.
func Query(
	preparer Preparer,
	query string,
	parameters []any,
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

func checkInsertResult(
	result sql.Result,
	err error,
	errorChecker ErrorChecker,
) (int64, error) {
	if err != nil {
		return 0, errorChecker.Check(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}

func queryMultiple[T Getter](
	preparer Preparer,
	query string,
	params []any,
	factoryFn func() T,
) ([]T, error) {
	rows, statement, err := Query(preparer, query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer statement.Close()

	return rowsToEntities[T](rows, factoryFn)
}

func querySingle[T Getter](
	preparer Preparer,
	query string,
	params []any,
	factoryFn func() T,
) (T, error) {
	var zero T

	stmt, err := preparer.Prepare(query)
	if err != nil {
		return zero, err
	}
	defer stmt.Close()

	return rowToEntity[T](stmt.QueryRow(params...), factoryFn)
}

func rowToEntity[T Getter](
	row Row,
	factoryFn func() T,
) (T, error) {
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

func rowsToEntities[T Getter](
	rows Rows,
	factoryFn func() T,
) ([]T, error) {
	var results []T
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

func checkUpdateResult(
	result sql.Result,
	err error,
	errorChecker ErrorChecker,
) (int64, error) {
	if err != nil {
		return 0, errorChecker.Check(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
