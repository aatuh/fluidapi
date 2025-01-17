package query

import (
	"fmt"
	"strings"

	util "github.com/pakkasys/fluidapi/database/util"
	"github.com/pakkasys/fluidapi/endpoint/page"
)

// GetOptions is the options struct used for get queries.
type GetOptions struct {
	Selectors   []util.Selector
	Orders      []util.Order
	Page        *page.Page
	Joins       []util.Join
	Projections []util.Projection
	Lock        bool
}

// Get returns a get query.
//
//   - tableName: The name of the database table.
//   - dbOptions: The options for the query.
func Get(tableName string, dbOptions *GetOptions) (string, []any) {
	whereClause, whereValues := whereClause(dbOptions.Selectors)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(
		"SELECT %s",
		strings.Join(projectionsToStrings(dbOptions.Projections), ","),
	))
	builder.WriteString(fmt.Sprintf(" FROM `%s`", tableName))
	if len(dbOptions.Joins) != 0 {
		builder.WriteString(" " + joinClause(dbOptions.Joins))
	}
	if whereClause != "" {
		builder.WriteString(" " + whereClause)
	}
	if len(dbOptions.Orders) != 0 {
		builder.WriteString(" " + getOrderClauseFromOrders(dbOptions.Orders))
	}
	if dbOptions.Page != nil {
		builder.WriteString(" " + getLimitOffsetClauseFromPage(dbOptions.Page))
	}
	if dbOptions.Lock {
		builder.WriteString(" FOR UPDATE")
	}

	return builder.String(), whereValues
}

func projectionsToStrings(projections []util.Projection) []string {
	if len(projections) == 0 {
		return []string{"*"}
	}

	projectionStrings := make([]string, len(projections))
	for i, projection := range projections {
		projectionStrings[i] = projection.String()
	}
	return projectionStrings
}

func joinClause(joins []util.Join) string {
	var joinClause string
	for _, join := range joins {
		if joinClause != "" {
			joinClause += " "
		}
		joinClause += fmt.Sprintf(
			"%s JOIN `%s` ON %s = %s",
			join.Type,
			join.Table,
			join.OnLeft.String(),
			join.OnRight.String(),
		)
	}
	return joinClause
}

func whereClause(selectors []util.Selector) (string, []any) {
	whereColumns, whereValues := processSelectors(selectors)

	var whereClause string
	if len(whereColumns) > 0 {
		whereClause = "WHERE " + strings.Join(whereColumns, " AND ")
	}

	return strings.Trim(whereClause, " "), whereValues
}

func getLimitOffsetClauseFromPage(page *page.Page) string {
	if page == nil {
		return ""
	}

	return fmt.Sprintf(
		"LIMIT %d OFFSET %d",
		page.Limit,
		page.Offset,
	)
}

func getOrderClauseFromOrders(orders []util.Order) string {
	if len(orders) == 0 {
		return ""
	}

	orderClause := "ORDER BY"
	for _, readOrder := range orders {
		if readOrder.Table == "" {
			orderClause += fmt.Sprintf(
				" `%s` %s,",
				readOrder.Field,
				readOrder.Direction,
			)
		} else {
			orderClause += fmt.Sprintf(
				" `%s`.`%s` %s,",
				readOrder.Table,
				readOrder.Field,
				readOrder.Direction,
			)
		}
	}

	return strings.TrimSuffix(orderClause, ",")
}
