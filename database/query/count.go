package query

import (
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/clause"
)

// CountOptions is the options struct used for count queries.
type CountOptions struct {
	Selectors []clause.Selector
	Joins     []clause.Join
}

// Count returns a count query.
//
//   - tableName: The name of the database table.
//   - dbOptions: The options for the query.
func Count(
	tableName string,
	dbOptions *CountOptions,
) (string, []any) {
	whereClause, whereValues := whereClause(dbOptions.Selectors)

	query := strings.Trim(fmt.Sprintf(
		"SELECT COUNT(*) FROM `%s` %s %s",
		tableName,
		joinClause(dbOptions.Joins),
		whereClause,
	), " ")

	return query, whereValues
}
