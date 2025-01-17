package query

import (
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/util"
)

// UpdateField is the options struct used for update queries.
type UpdateField struct {
	Field string
	Value any
}

// UpdateQuery returns the SQL query and values for an update query.
//
//   - tableName: The name of the database table.
//   - updateFields: The fields to update.
//   - selectors: The selectors for the entities to update.
func UpdateQuery(
	tableName string,
	updateFields []UpdateField,
	selectors []util.Selector,
) (string, []any) {
	whereColumns, whereValues := processSelectors(selectors)

	setClause, values := getSetClause(updateFields)
	values = append(values, whereValues...)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(
		"UPDATE `%s` SET %s",
		tableName,
		setClause,
	))
	if len(whereColumns) != 0 {
		builder.WriteString(" " + getWhereClause(whereColumns))
	}

	return builder.String(), values
}

func getWhereClause(whereColumns []string) string {
	whereClause := ""
	if len(whereColumns) > 0 {
		whereClause = "WHERE " + strings.Join(whereColumns, " AND ")
	}
	return whereClause
}

func getSetClause(updates []UpdateField) (string, []any) {
	setClauseParts := make([]string, len(updates))
	values := make([]any, len(updates))

	for i, update := range updates {
		setClauseParts[i] = fmt.Sprintf(
			"%s = ?",
			update.Field,
		)
		values[i] = update.Value
	}

	return strings.Join(setClauseParts, ", "), values
}
