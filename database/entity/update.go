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
	query, values := query.UpdateQuery(tableName, updateFields, selectors)
	result, err := Exec(preparer, query, values)
	return checkUpdateResult(result, err, errorChecker)
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
