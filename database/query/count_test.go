package query

import (
	"testing"

	"github.com/pakkasys/fluidapi/database/clause"
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
		Selectors: []clause.Selector{
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
		Joins: []clause.Join{
			{
				Type:  clause.JoinTypeInner,
				Table: "other_table",
				OnLeft: clause.ColumnSelector{
					Table:   "test_table",
					Columnn: "id",
				},
				OnRight: clause.ColumnSelector{
					Table:   "other_table",
					Columnn: "ref_id",
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
		Selectors: []clause.Selector{
			{Table: "test_table", Field: "id", Predicate: "=", Value: 1},
		},
		Joins: []clause.Join{
			{
				Type:  clause.JoinTypeInner,
				Table: "other_table",
				OnLeft: clause.ColumnSelector{
					Table:   "test_table",
					Columnn: "id",
				},
				OnRight: clause.ColumnSelector{
					Table:   "other_table",
					Columnn: "ref_id",
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
