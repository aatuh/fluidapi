package entity

import (
	"database/sql"

	"github.com/pakkasys/fluidapi/database/query"
	"github.com/pakkasys/fluidapi/database/util"
)

// Create creates an entity in the database.
//
//   - entity: The entity to insert.
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - inserter: The function used to get the columns and values to insert.
func Create[T any](
	entity *T,
	preparer util.Preparer,
	tableName string,
	inserter query.InsertedValues[*T],
	errorChecker ErrorChecker,
) (int64, error) {
	res, err := insert(preparer, entity, tableName, inserter)
	return checkInsertResult(res, err, errorChecker)
}

// CreateMany creates entities in the database.
//
//   - entities: The entities to insert.
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - inserter: The function used to get the columns and values to insert.
func CreateMany[T any](
	entities []*T,
	preparer util.Preparer,
	tableName string,
	inserter query.InsertedValues[*T],
	errorChecker ErrorChecker,
) (int64, error) {
	if len(entities) == 0 {
		return 0, nil
	}
	res, err := insertMany(
		preparer,
		entities,
		tableName,
		inserter,
	)
	return checkInsertResult(res, err, errorChecker)
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

func insert[T any](
	preparer util.Preparer,
	entity *T,
	tableName string,
	inserter query.InsertedValues[*T],
) (sql.Result, error) {
	query, values := query.Insert(entity, tableName, inserter)

	statement, err := preparer.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.Exec(values...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func insertMany[T any](
	preparer util.Preparer,
	entities []*T,
	tableName string,
	inserter query.InsertedValues[*T],
) (sql.Result, error) {
	query, values := query.InsertMany(entities, tableName, inserter)

	statement, err := preparer.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.Exec(values...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
