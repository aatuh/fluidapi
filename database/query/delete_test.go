package query

import (
	"strings"
	"testing"

	"github.com/pakkasys/fluidapi/database/clause"
	"github.com/stretchr/testify/assert"
)

// TestWriteDeleteOptions_WithLimitAndOrders tests writeDeleteOptions with both
// limit and orders.
func TestWriteDeleteOptions_WithLimitAndOrders(t *testing.T) {
	// Create a DeleteOptions with a limit and orders
	orders := []clause.Order{
		{Table: "user", Field: "name", Direction: "ASC"},
		{Table: "user", Field: "age", Direction: "DESC"},
	}
	opts := DeleteOptions{Limit: 10, Orders: orders}

	// Create a string builder for the SQL query
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM `user` WHERE id = 1")

	writeDeleteOptions(&builder, &opts)

	expectedSQL := "DELETE FROM `user` WHERE id = 1 ORDER BY `user`.`name` ASC, `user`.`age` DESC LIMIT 10"

	assert.Equal(t, expectedSQL, builder.String())
}

// TestWriteDeleteOptions_WithOnlyOrders tests writeDeleteOptions with only
// orders and no limit.
func TestWriteDeleteOptions_WithOnlyOrders(t *testing.T) {
	// Create a DeleteOptions with only orders
	orders := []clause.Order{
		{Table: "user", Field: "name", Direction: "ASC"},
	}
	opts := DeleteOptions{Limit: 0, Orders: orders}

	// Create a string builder for the SQL query
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM `user` WHERE id = 1")

	writeDeleteOptions(&builder, &opts)

	expectedSQL := "DELETE FROM `user` WHERE id = 1 ORDER BY `user`.`name` ASC"

	assert.Equal(t, expectedSQL, builder.String())
}

// TestWriteDeleteOptions_WithOnlyLimit tests writeDeleteOptions with only a
// limit and no orders.
func TestWriteDeleteOptions_WithOnlyLimit(t *testing.T) {
	// Create a DeleteOptions with only a limit
	opts := DeleteOptions{Limit: 5, Orders: nil}

	// Create a string builder for the SQL query
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM `user` WHERE id = 1")

	writeDeleteOptions(&builder, &opts)

	expectedSQL := "DELETE FROM `user` WHERE id = 1 LIMIT 5"

	assert.Equal(t, expectedSQL, builder.String())
}

// TestWriteDeleteOptions_WithNoOptions tests writeDeleteOptions with no limit
// and no orders.
func TestWriteDeleteOptions_WithNoOptions(t *testing.T) {
	// Create an empty DeleteOptions with no limit and no orders
	opts := DeleteOptions{Limit: 0, Orders: nil}

	// Create a string builder for the SQL query
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM `user` WHERE id = 1")

	writeDeleteOptions(&builder, &opts)

	expectedSQL := "DELETE FROM `user` WHERE id = 1"

	assert.Equal(t, expectedSQL, builder.String())
}
