package entity

import (
	"database/sql"

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
	result, err := delete(preparer, tableName, selectors, opts)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func delete(
	preparer util.Preparer,
	tableName string,
	selectors []util.Selector,
	opts *query.DeleteOptions,
) (sql.Result, error) {
	query, whereValues := query.Delete(tableName, selectors, opts)

	statement, err := preparer.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.Exec(whereValues...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
