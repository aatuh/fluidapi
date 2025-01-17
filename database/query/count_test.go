package query

import (
	"testing"

	"github.com/pakkasys/fluidapi/database/util"
	"github.com/stretchr/testify/assert"
)

// TestBuildBaseCountQuery_NoSelectorsNoJoins tests BuildBaseCountQuery with no
// selectors or joins.
func TestBuildBaseCountQuery_NoSelectorsNoJoins(t *testing.T) {
	tableName := "test_table"
	dbOptions := &CountOptions{}

	query, whereValues := Count(tableName, dbOptions)

	expectedQuery := "SELECT COUNT(*) FROM `test_table`"
	expectedValues := []any{}

	assert.Equal(t, expectedQuery, query)
	assert.ElementsMatch(t, expectedValues, whereValues)
}

// TestBuildBaseCountQuery_WithSelectors tests BuildBaseCountQuery only
// selectors.
func TestBuildBaseCountQuery_WithSelectors(t *testing.T) {
	tableName := "test_table"
	dbOptions := &CountOptions{
		Selectors: []util.Selector{
			{Table: "test_table", Field: "id", Predicate: "=", Value: 1},
		},
	}

	query, whereValues := Count(tableName, dbOptions)

	expectedQuery := "SELECT COUNT(*) FROM `test_table`  WHERE `test_table`.`id` = ?"
	expectedValues := []any{1}

	assert.Equal(t, expectedQuery, query)
	assert.ElementsMatch(t, expectedValues, whereValues)
}

// TestBuildBaseCountQuery_WithJoins tests BuildBaseCountQuery with joins only.
func TestBuildBaseCountQuery_WithJoins(t *testing.T) {
	tableName := "test_table"
	dbOptions := &CountOptions{
		Joins: []util.Join{
			{
				Type:  util.JoinTypeInner,
				Table: "other_table",
				OnLeft: util.ColumSelector{
					Table:  "test_table",
					Column: "id",
				},
				OnRight: util.ColumSelector{
					Table:  "other_table",
					Column: "ref_id",
				},
			},
		},
	}
	query, whereValues := Count(tableName, dbOptions)

	expectedQuery := "SELECT COUNT(*) FROM `test_table` INNER JOIN `other_table` ON `test_table`.`id` = `other_table`.`ref_id`"
	expectedValues := []any{}

	assert.Equal(t, expectedQuery, query)
	assert.ElementsMatch(t, expectedValues, whereValues)
}

// TestBuildBaseCountQuery_WithSelectorsAndJoins tests BuildBaseCountQuery with
// both selectors and joins.
func TestBuildBaseCountQuery_WithSelectorsAndJoins(t *testing.T) {
	tableName := "test_table"
	dbOptions := &CountOptions{
		Selectors: []util.Selector{
			{Table: "test_table", Field: "id", Predicate: "=", Value: 1},
		},
		Joins: []util.Join{
			{
				Type:  util.JoinTypeInner,
				Table: "other_table",
				OnLeft: util.ColumSelector{
					Table:  "test_table",
					Column: "id",
				},
				OnRight: util.ColumSelector{
					Table:  "other_table",
					Column: "ref_id",
				},
			},
		},
	}

	query, whereValues := Count(tableName, dbOptions)

	expectedQuery := "SELECT COUNT(*) FROM `test_table` INNER JOIN `other_table` ON `test_table`.`id` = `other_table`.`ref_id` WHERE `test_table`.`id` = ?"
	expectedValues := []any{1}

	assert.Equal(t, expectedQuery, query)
	assert.ElementsMatch(t, expectedValues, whereValues)
}
