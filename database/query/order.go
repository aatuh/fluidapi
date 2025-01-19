package query

import (
	"fmt"
	"strings"

	"github.com/pakkasys/fluidapi/database/clause"
)

// GetOrderClauseFromOrders returns an order clause from the provided orders.
//
//   - orders: The orders to get the clause for.
func GetOrderClauseFromOrders(orders []clause.Order) string {
	if len(orders) == 0 {
		return ""
	}

	orderClause := "ORDER BY"
	for _, order := range orders {
		if order.Table == "" {
			orderClause += fmt.Sprintf(
				" `%s` %s,",
				order.Field,
				order.Direction,
			)
		} else {
			orderClause += fmt.Sprintf(
				" `%s`.`%s` %s,",
				order.Table,
				order.Field,
				order.Direction,
			)
		}
	}

	return strings.TrimSuffix(orderClause, ",")
}
