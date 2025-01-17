package entity

import (
	"github.com/pakkasys/fluidapi/database/query"
	"github.com/pakkasys/fluidapi/database/util"
)

// Delete deletes entities from the database.
//
//   - db: The database connection.
//   - tableName: The name of the database table.
//   - selectors: The selectors for the entities to delete.
//   - opts: The options for the query.
func Delete(
	preparer util.Preparer,
	tableName string,
	selectors []util.Selector,
	opts *query.DeleteOptions,
) (int64, error) {
	query, whereValues := query.Delete(tableName, selectors, opts)

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
