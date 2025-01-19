package query

import (
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/clause"
	"github.com/pakkasys/fluidapi/endpoint/page"
)

// GetOptions is the options struct used for get queries.
type GetOptions struct {
	Selectors   []clause.Selector
	Orders      []clause.Order
	Page        *page.Page
	Joins       []clause.Join
	Projections []clause.Projection
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
		builder.WriteString(" " + GetOrderClauseFromOrders(dbOptions.Orders))
	}
	if dbOptions.Page != nil {
		builder.WriteString(" " + page.GetLimitOffsetClauseFromPage(dbOptions.Page))
	}
	if dbOptions.Lock {
		builder.WriteString(" FOR UPDATE")
	}

	return builder.String(), whereValues
}

func projectionsToStrings(projections []clause.Projection) []string {
	if len(projections) == 0 {
		return []string{"*"}
	}

	projectionStrings := make([]string, len(projections))
	for i, projection := range projections {
		projectionStrings[i] = ProjectionToString(projection)
	}
	return projectionStrings
}

func joinClause(joins []clause.Join) string {
	var joinClause string
	for _, join := range joins {
		if joinClause != "" {
			joinClause += " "
		}
		joinClause += fmt.Sprintf(
			"%s JOIN `%s` ON %s = %s",
			join.Type,
			join.Table,
			ColumnSelectorToString(join.OnLeft),
			ColumnSelectorToString(join.OnRight),
		)
	}
	return joinClause
}

func whereClause(selectors []clause.Selector) (string, []any) {
	whereColumnns, whereValues := ProcessSelectors(selectors)

	var whereClause string
	if len(whereColumnns) > 0 {
		whereClause = "WHERE " + strings.Join(whereColumnns, " AND ")
	}

	return strings.Trim(whereClause, " "), whereValues
}
