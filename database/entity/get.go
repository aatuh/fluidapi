package entity

import (
	"database/sql"
	"fmt"

	"github.com/pakkasys/fluidapi/database/query"
	util "github.com/pakkasys/fluidapi/database/util"
)

type RowScanner[T any] func(row util.Row, entity *T) error
type RowScannerMultiple[T any] func(rows util.Rows, entity *T) error

// Get returns a single entity.
//
//   - tableName: The table name.
//   - rowScanner: The function used to scan the row.
//   - preparer: The database preparer.
//   - dbOptions: The options for the query.
func Get[T any](
	tableName string,
	rowScanner RowScanner[T],
	preparer util.Preparer,
	dbOptions *query.GetOptions,
) (*T, error) {
	query, whereValues := query.Get(tableName, dbOptions)

	entity, err := querySingle(
		preparer,
		query,
		whereValues,
		rowScanner,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return entity, nil
}

// GetMany returns multiple entities.
//
//   - tableName: The table name.
//   - rowScannerMultiple: The function used to scan multiple rows.
//   - preparer: The database preparer.
//   - dbOptions: The options for the query.
func GetMany[T any](
	tableName string,
	rowScannerMultiple RowScannerMultiple[T],
	preparer util.Preparer,
	dbOptions *query.GetOptions,
) ([]T, error) {
	query, whereValues := query.Get(tableName, dbOptions)

	entities, err := queryMultiple(preparer, query, whereValues, rowScannerMultiple)
	if err != nil {
		if err == sql.ErrNoRows {
			return []T{}, nil
		}
		return nil, err
	}
	return entities, nil
}

func queryMultiple[T any](
	preparer util.Preparer,
	query string,
	params []any,
	rowScannerMultiple RowScannerMultiple[T],
) ([]T, error) {
	rows, statement, err := Query(preparer, query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer statement.Close()

	entities, err := rowsToEntities(rows, rowScannerMultiple)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func querySingle[T any](
	preparer util.Preparer,
	query string,
	params []any,
	rowScanner RowScanner[T],
) (*T, error) {
	statement, err := preparer.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var entity T

	row := statement.QueryRow(params...)
	err = rowScanner(row, &entity)
	if err != nil {
		return nil, err
	}
	if err := row.Err(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func rowsToEntities[T any](
	rows util.Rows,
	rowScannerMultiple RowScannerMultiple[T],
) ([]T, error) {
	if rowScannerMultiple == nil {
		return nil, fmt.Errorf("must provide rowScannerMultiple")
	}

	var entities []T
	for rows.Next() {
		var entity T
		err := rowScannerMultiple(rows, &entity)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}
