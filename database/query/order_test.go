package query

import (
	"testing"

	"github.com/pakkasys/fluidapi/database/clause"
	"github.com/stretchr/testify/assert"
)

// TestGetOrderClauseFromOrders_NoOrders tests the case where no orders are
// provided.
func TestGetOrderClauseFromOrders_NoOrders(t *testing.T) {
	orders := []clause.Order{}
	orderClause := GetOrderClauseFromOrders(orders)
	assert.Equal(t, "", orderClause)
}

// TestGetOrderClauseFromOrders_WithoutTable tests the case where there is no
// table in the order.
func TestGetOrderClauseFromOrders_WithoutTable(t *testing.T) {
	orders := []clause.Order{
		{Field: "name", Direction: "ASC"},
	}

	orderClause := GetOrderClauseFromOrders(orders)

	expected := "ORDER BY `name` ASC"
	assert.Equal(t, expected, orderClause)
}

// TestGetOrderClauseFromOrders_SingleOrder tests the case where a single order
// is provided.
func TestGetOrderClauseFromOrders_SingleOrder(t *testing.T) {
	orders := []clause.Order{
		{Table: "user", Field: "name", Direction: "ASC"},
	}

	orderClause := GetOrderClauseFromOrders(orders)

	expected := "ORDER BY `user`.`name` ASC"
	assert.Equal(t, expected, orderClause)
}

// TestGetOrderClauseFromOrders_MultipleOrders tests the case where multiple
// orders are provided.
func TestGetOrderClauseFromOrders_MultipleOrders(t *testing.T) {
	orders := []clause.Order{
		{Table: "user", Field: "name", Direction: "ASC"},
		{Table: "user", Field: "age", Direction: "DESC"},
	}

	orderClause := GetOrderClauseFromOrders(orders)

	expected := "ORDER BY `user`.`name` ASC, `user`.`age` DESC"
	assert.Equal(t, expected, orderClause)
}

// TestGetOrderClauseFromOrders_EmptyFields tests the case where orders have
// empty fields.
func TestGetOrderClauseFromOrders_EmptyFields(t *testing.T) {
	orders := []clause.Order{
		{Table: "", Field: "", Direction: "ASC"},
	}

	orderClause := GetOrderClauseFromOrders(orders)

	// Expect an ORDER BY clause with empty table and field
	expected := "ORDER BY `` ASC"
	assert.Equal(t, expected, orderClause)
}
