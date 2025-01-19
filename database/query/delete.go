package query

import (
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/clause"
)

// DeleteOptions is the options struct used for delete queries.
type DeleteOptions struct {
	Limit  int
	Orders []clause.Order
}

// Delete returns the SQL query string and the values for the query.
//
//   - tableName: The name of the database table.
//   - selectors: The selectors for the entities to delete.
//   - opts: The options for the query.
func Delete(
	tableName string,
	selectors []clause.Selector,
	opts *DeleteOptions,
) (string, []any) {
	whereColumns, whereValues := ProcessSelectors(selectors)

	whereClause := ""
	if len(whereColumns) > 0 {
		whereClause = "WHERE " + strings.Join(whereColumns, " AND ")
	}

	builder := strings.Builder{}
	builder.WriteString(
		fmt.Sprintf("DELETE FROM `%s` %s", tableName, whereClause),
	)

	if opts != nil {
		writeDeleteOptions(&builder, opts)
	}

	return builder.String(), whereValues
}

func writeDeleteOptions(
	builder *strings.Builder,
	opts *DeleteOptions,
) {
	orderClause := GetOrderClauseFromOrders(opts.Orders)
	if orderClause != "" {
		builder.WriteString(" " + orderClause)
	}

	limit := opts.Limit
	if limit > 0 {
		builder.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	}
}
