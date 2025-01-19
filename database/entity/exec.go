package entity

import (
	"github.com/pakkasys/fluidapi/database"
)

// Exec runs a query and returns the result.
//
//   - preparer: The preparer used to prepare the query.
//   - query: The query string.
//   - parameters: The parameters for the query.
func Exec(
	preparer database.Preparer,
	query string,
	parameters []any,
) (database.Result, error) {
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

// Query runs a query and returns the rows. The returned rows and stmt objects
// must be managed manually by the caller after successful query execution.
//
//   - preparer: The preparer used to prepare the query.
//   - query: The query string.
//   - parameters: The parameters for the query.
func Query(
	preparer database.Preparer,
	query string,
	parameters []any,
) (database.Rows, database.Stmt, error) {
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
