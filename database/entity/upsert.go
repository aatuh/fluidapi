package entity

import (
	"fmt"

	"github.com/pakkasys/fluidapi/database/query"
	"github.com/pakkasys/fluidapi/database/util"
)

// Upsert upserts an entity.
//
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - entity: The entity to upsert.
//   - insertedValues: The function used to get the columns and values to insert.
//   - updateProjections: The projections of the entity to update.
func Upsert[T any](
	preparer util.Preparer,
	tableName string,
	entity *T,
	insertedValues query.InsertedValues[*T],
	updateProjections []util.Projection,
	errorChecker ErrorChecker,
) (int64, error) {
	return UpsertMany(
		preparer,
		tableName,
		[]*T{entity},
		insertedValues,
		updateProjections,
		errorChecker,
	)
}

// UpsertMany upserts multiple entities.
//
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - entities: The entities to upsert.
//   - insertedValues: The function used to get the columns and values to insert.
//   - updateProjections: The projections of the entities to update.
func UpsertMany[T any](
	preparer util.Preparer,
	tableName string,
	entities []*T,
	insertedValues query.InsertedValues[*T],
	updateProjections []util.Projection,
	errorChecker ErrorChecker,
) (int64, error) {
	if len(entities) == 0 {
		return 0, fmt.Errorf("must provide entities to upsert")
	}
	if len(updateProjections) == 0 {
		return 0, fmt.Errorf("must provide update projections")
	}
	if len(updateProjections[0].Alias) == 0 {
		return 0, fmt.Errorf("must provide update projections alias")
	}

	query, values := query.UpsertMany(
		entities,
		tableName,
		insertedValues,
		updateProjections,
	)
	statement, err := preparer.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(values...)
	if err != nil {
		return 0, err
	}

	return checkInsertResult(result, err, errorChecker)
}
