package query

import (
	"testing"

	"github.com/pakkasys/fluidapi/database/clause"
	"github.com/stretchr/testify/assert"
)

// TestUpdateQuery_SingleUpdate tests the case where a single update is
// provided.
func TestUpdateQuery_SingleUpdate(t *testing.T) {
	updates := []UpdateField{
		{Field: "name", Value: "Alice"},
	}
	selectors := []clause.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	query, values := UpdateQuery("user", updates, selectors)

	expectedQuery := "UPDATE `user` SET name = ? WHERE `user`.`id` = ?"
	expectedValues := []any{"Alice", 1}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpdateQuery_MultipleUpdates tests the case where multiple updates are
// provided.
func TestUpdateQuery_MultipleUpdates(t *testing.T) {
	updates := []UpdateField{
		{Field: "name", Value: "Alice"},
		{Field: "age", Value: 30},
	}
	selectors := []clause.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	query, values := UpdateQuery("user", updates, selectors)

	expectedQuery := "UPDATE `user` SET name = ?, age = ? WHERE `user`.`id` = ?"
	expectedValues := []any{"Alice", 30, 1}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpdateQuery_NoUpdates tests the case where no updates are provided.
func TestUpdateQuery_NoUpdates(t *testing.T) {
	updates := []UpdateField{}
	selectors := []clause.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	query, values := UpdateQuery("user", updates, selectors)

	expectedQuery := "UPDATE `user` SET  WHERE `user`.`id` = ?"
	expectedValues := []any{1}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpdateQuery_NoSelectors tests the case where no selectors are provided.
func TestUpdateQuery_NoSelectors(t *testing.T) {
	updates := []UpdateField{
		{Field: "name", Value: "Alice"},
	}

	selectors := []clause.Selector{} // No selectors

	query, values := UpdateQuery("user", updates, selectors)

	expectedQuery := "UPDATE `user` SET name = ?"
	expectedValues := []any{"Alice"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestUpdateQuery_EmptyFields tests the case where updates and selectors have
// empty fields.
func TestUpdateQuery_EmptyFields(t *testing.T) {
	updates := []UpdateField{
		{Field: "", Value: "Unknown"},
	}
	selectors := []clause.Selector{
		{Table: "", Field: "", Predicate: "=", Value: nil},
	}

	query, values := UpdateQuery("user", updates, selectors)

	expectedQuery := "UPDATE `user` SET  = ? WHERE `` IS NULL"
	expectedValues := []any{"Unknown"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestGetWhereClause_NoConditions tests the case where no conditions are
// provided.
func TestGetWhereClause_NoConditions(t *testing.T) {
	whereColumns := []string{}

	whereClause := getWhereClause(whereColumns)

	expectedWhereClause := ""
	assert.Equal(t, expectedWhereClause, whereClause)
}

// TestGetWhereClause_SingleCondition tests the case where a single condition is
// provided.
func TestGetWhereClause_SingleCondition(t *testing.T) {
	whereColumns := []string{"`user`.`id` = ?"}

	whereClause := getWhereClause(whereColumns)

	expectedWhereClause := "WHERE `user`.`id` = ?"
	assert.Equal(t, expectedWhereClause, whereClause)
}

// TestGetWhereClause_MultipleConditions tests the case where multiple
// conditions are provided.
func TestGetWhereClause_MultipleConditions(t *testing.T) {
	whereColumns := []string{"`user`.`id` = ?", "`user`.`age` > 18"}

	whereClause := getWhereClause(whereColumns)

	expectedWhereClause := "WHERE `user`.`id` = ? AND `user`.`age` > 18"
	assert.Equal(t, expectedWhereClause, whereClause)
}

// TestGetSetClause_SingleUpdate tests the case where a single update is
// provided.
func TestGetSetClause_SingleUpdate(t *testing.T) {
	updates := []UpdateField{
		{Field: "name", Value: "Alice"},
	}

	setClause, values := getSetClause(updates)

	expectedSetClause := "name = ?"
	expectedValues := []any{"Alice"}

	assert.Equal(t, expectedSetClause, setClause)
	assert.Equal(t, expectedValues, values)
}

// TestGetSetClause_MultipleUpdates tests the case where multiple updates are
// provided.
func TestGetSetClause_MultipleUpdates(t *testing.T) {
	updates := []UpdateField{
		{Field: "name", Value: "Alice"},
		{Field: "age", Value: 30},
	}

	setClause, values := getSetClause(updates)

	expectedSetClause := "name = ?, age = ?"
	expectedValues := []any{"Alice", 30}

	assert.Equal(t, expectedSetClause, setClause)
	assert.Equal(t, expectedValues, values)
}

// TestGetSetClause_NoUpdates tests the case where no updates are provided.
func TestGetSetClause_NoUpdates(t *testing.T) {
	updates := []UpdateField{}

	setClause, values := getSetClause(updates)

	expectedSetClause := ""
	expectedValues := []any{}

	assert.Equal(t, expectedSetClause, setClause)
	assert.Equal(t, expectedValues, values)
}

// TestGetSetClause_EmptyField tests the case where an update has an empty
// field.
func TestGetSetClause_EmptyField(t *testing.T) {
	updates := []UpdateField{
		{Field: "", Value: "Unknown"},
	}

	setClause, values := getSetClause(updates)

	expectedSetClause := " = ?"
	expectedValues := []any{"Unknown"}

	assert.Equal(t, expectedSetClause, setClause)
	assert.Equal(t, expectedValues, values)
}
