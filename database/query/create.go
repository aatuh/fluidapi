package query

import (
	"fmt"
	"strings"
)

// InsertedValues is a function used to get the columns and values to insert.
type InsertedValues[T any] func(entity T) (columns []string, values []any)

// Insert returns the query and values to insert an entity.
//
//   - entity: The entity to insert.
//   - tableName: The name of the database table.
//   - insertedValues: Function used to get the columns and values to insert.
func Insert[T any](
	entity *T,
	tableName string,
	insertedValues InsertedValues[*T],
) (string, []any) {
	columns, values := insertedValues(entity)
	columnNames := getInsertQueryColumnNames(columns)

	valuePlaceholders := strings.TrimSuffix(
		strings.Repeat("?, ", len(values)),
		", ",
	)

	query := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s)",
		tableName,
		columnNames,
		valuePlaceholders,
	)

	return query, values
}

// InsertMany returns the query and values to insert multiple entities.
//
//   - entities: The entities to insert.
//   - tableName: The name of the database table.
//   - insertedValues: Function used to get the columns and values to insert.
func InsertMany[T any](
	entities []*T,
	tableName string,
	insertedValues InsertedValues[*T],
) (string, []any) {
	if len(entities) == 0 {
		return "", nil
	}

	columns, _ := insertedValues(entities[0])
	columnNames := getInsertQueryColumnNames(columns)

	var allValues []any
	valuePlaceholders := make([]string, len(entities))
	for i, entity := range entities {
		_, values := insertedValues(entity)
		placeholders := make([]string, len(values))
		for j := range values {
			placeholders[j] = "?"
		}
		valuePlaceholders[i] = "(" + strings.Join(placeholders, ", ") + ")"
		allValues = append(allValues, values...)
	}

	query := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES %s",
		tableName,
		columnNames,
		strings.Join(valuePlaceholders, ", "),
	)

	return query, allValues
}

func getInsertQueryColumnNames(columns []string) string {
	wrappedColumns := make([]string, len(columns))
	for i, column := range columns {
		wrappedColumns[i] = "`" + column + "`"
	}
	columnNames := strings.Join(wrappedColumns, ", ")
	return columnNames
}
