package entity

import (
	"database/sql"

	"github.com/pakkasys/fluidapi/database/query"
	"github.com/pakkasys/fluidapi/database/util"
)

// Update updates entities in the database.
//
//   - db: The database connection to use.
//   - tableName: The name of the database table.
//   - selectors: The selectors of the entities to update.
//   - updates: The updates to apply to the entities.
func Update(
	preparer util.Preparer,
	tableName string,
	selectors []util.Selector,
	updateFields []query.UpdateField,
	errorChecker ErrorChecker,
) (int64, error) {
	if len(updateFields) == 0 {
		return 0, nil
	}
	res, err := update(preparer, tableName, updateFields, selectors)
	return checkUpdateResult(res, err, errorChecker)
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

func update(
	preparer util.Preparer,
	tableName string,
	updateFields []query.UpdateField,
	selectors []util.Selector,
) (sql.Result, error) {
	query, values := query.UpdateQuery(tableName, updateFields, selectors)

	statement, err := preparer.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.Exec(values...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
