package query

import (
	"testing"

	"github.com/pakkasys/fluidapi/database/clause"
	"github.com/stretchr/testify/assert"
)

// TestColumnSelectorString_NormalCase tests the normal operation of the String
// method.
func TestColumnSelectorString_NormalCase(t *testing.T) {
	selector := clause.ColumnSelector{
		Table:   "users",
		Columnn: "id",
	}

	result := ColumnSelectorToString(selector)
	expected := "`users`.`id`"

	assert.Equal(t, expected, result)
}

// TestColumnSelectorString_EmptyTable tests the case where the Table is empty.
func TestColumnSelectorString_EmptyTable(t *testing.T) {
	selector := clause.ColumnSelector{
		Table:   "",
		Columnn: "id",
	}

	result := ColumnSelectorToString(selector)
	expected := "``.`id`"

	assert.Equal(t, expected, result)
}

// TestColumnSelectorString_EmptyColumnn tests the case where the Columnn is empty.
func TestColumnSelectorString_EmptyColumnn(t *testing.T) {
	selector := clause.ColumnSelector{
		Table:   "users",
		Columnn: "",
	}

	result := ColumnSelectorToString(selector)
	expected := "`users`.``"

	assert.Equal(t, expected, result)
}

// TestColumnSelectorString_EmptyTableAndColumnn tests the case where both the
// Table and Columnn are empty.
func TestColumnSelectorString_EmptyTableAndColumnn(t *testing.T) {
	selector := clause.ColumnSelector{
		Table:   "",
		Columnn: "",
	}

	result := ColumnSelectorToString(selector)
	expected := "``.``"

	assert.Equal(t, expected, result)
}
