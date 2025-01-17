package entity

import (
	"database/sql"

	"github.com/pakkasys/fluidapi/database/query"
	"github.com/pakkasys/fluidapi/database/util"
)

// Insert creates an entity in the database.
//
//   - entity: The entity to insert.
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - inserter: The function used to get the columns and values to insert.
func Insert[T any](
	entity *T,
	preparer util.Preparer,
	tableName string,
	inserter query.InsertedValues[*T],
	errorChecker ErrorChecker,
) (int64, error) {
	query, values := query.Insert(entity, tableName, inserter)
	result, err := Exec(preparer, query, values)
	return checkInsertResult(result, err, errorChecker)
}

// InsertMany creates many entities in the database.
//
//   - entities: The entities to insert.
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - inserter: The function used to get the columns and values to insert.
func InsertMany[T any](
	entities []*T,
	preparer util.Preparer,
	tableName string,
	inserter query.InsertedValues[*T],
	errorChecker ErrorChecker,
) (int64, error) {
	if len(entities) == 0 {
		return 0, nil
	}
	query, values := query.InsertMany(entities, tableName, inserter)
	result, err := Exec(preparer, query, values)
	return checkInsertResult(result, err, errorChecker)
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
