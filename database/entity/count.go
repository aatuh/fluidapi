package entity

import (
	"github.com/pakkasys/fluidapi/database"
	"github.com/pakkasys/fluidapi/database/query"
)

// Count counts the number of entities in the database.
//
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - dbOptions: The options for the query.
func Count(
	preparer database.Preparer,
	tableName string,
	dbOptions *query.CountOptions,
) (int, error) {
	query, whereValues := query.Count(tableName, dbOptions)

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
