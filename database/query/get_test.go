package query

import (
	"testing"

	"github.com/pakkasys/fluidapi/database/clause"
	"github.com/pakkasys/fluidapi/endpoint/page"
	"github.com/stretchr/testify/assert"
)

// TestProjectionsToStrings_NoProjections tests the case where no projections
// are provided.
func TestProjectionsToStrings_NoProjections(t *testing.T) {
	projections := []clause.Projection{}
	projectionStrings := projectionsToStrings(projections)
	assert.Equal(t, []string{"*"}, projectionStrings)
}

// TestProjectionsToStrings_SingleProjection tests the case where a single
// projection is provided.
func TestProjectionsToStrings_SingleProjection(t *testing.T) {
	projections := []clause.Projection{
		{Table: "user", Column: "name"},
	}

	projectionStrings := projectionsToStrings(projections)

	expected := []string{"`user`.`name`"}
	assert.Equal(t, expected, projectionStrings)
}

// TestProjectionsToStrings_MultipleProjections tests the case where multiple
// projections are provided.
func TestProjectionsToStrings_MultipleProjections(t *testing.T) {
	projections := []clause.Projection{
		{Table: "user", Column: "name"},
		{Table: "orders", Column: "order_id"},
	}

	projectionStrings := projectionsToStrings(projections)

	expected := []string{"`user`.`name`", "`orders`.`order_id`"}
	assert.Equal(t, expected, projectionStrings)
}

// TestProjectionsToStrings_EmptyFields tests the case where a projection has
// empty fields.
func TestProjectionsToStrings_EmptyFields(t *testing.T) {
	projections := []clause.Projection{
		{Table: "", Column: ""},
	}

	projectionStrings := projectionsToStrings(projections)

	expected := []string{"``"}
	assert.Equal(t, expected, projectionStrings)
}

// TestJoinClause_NoJoins tests the case where no joins are provided.
func TestJoinClause_NoJoins(t *testing.T) {
	joins := []clause.Join{}
	joinClause := joinClause(joins)
	assert.Equal(t, "", joinClause)
}

// TestJoinClause_SingleJoin tests the case where a single join is provided.
func TestJoinClause_SingleJoin(t *testing.T) {
	joins := []clause.Join{
		{
			Type:  clause.JoinTypeInner,
			Table: "orders",
			OnLeft: clause.ColumnSelector{
				Table:   "user",
				Columnn: "id",
			},
			OnRight: clause.ColumnSelector{
				Table:   "orders",
				Columnn: "user_id",
			},
		},
	}

	joinClause := joinClause(joins)

	expected := "INNER JOIN `orders` ON `user`.`id` = `orders`.`user_id`"
	assert.Equal(t, expected, joinClause)
}

// TestJoinClause_MultipleJoins tests the case where multiple joins are
// provided.
func TestJoinClause_MultipleJoins(t *testing.T) {
	joins := []clause.Join{
		{
			Type:  clause.JoinTypeInner,
			Table: "order",
			OnLeft: clause.ColumnSelector{
				Table:   "user",
				Columnn: "id",
			},
			OnRight: clause.ColumnSelector{
				Table:   "order",
				Columnn: "user_id",
			},
		},
		{
			Type:  clause.JoinTypeLeft,
			Table: "payments",
			OnLeft: clause.ColumnSelector{
				Table:   "user",
				Columnn: "id",
			},
			OnRight: clause.ColumnSelector{
				Table:   "payments",
				Columnn: "user_id",
			},
		},
	}

	joinClause := joinClause(joins)

	// Expect multiple JOIN clauses
	expected := "INNER JOIN `order` ON `user`.`id` = `order`.`user_id` LEFT JOIN `payments` ON `user`.`id` = `payments`.`user_id`"
	assert.Equal(t, expected, joinClause)
}

// TestJoinClause_EmptyFields tests the case where a join has empty fields.
func TestJoinClause_EmptyFields(t *testing.T) {
	joins := []clause.Join{
		{
			Type:  clause.JoinTypeInner,
			Table: "",
			OnLeft: clause.ColumnSelector{
				Table:   "",
				Columnn: "",
			},
			OnRight: clause.ColumnSelector{
				Table:   "",
				Columnn: "",
			},
		},
	}

	joinClause := joinClause(joins)

	// Expect a malformed JOIN clause with empty fields
	expected := "INNER JOIN `` ON ``.`` = ``.``"
	assert.Equal(t, expected, joinClause)
}

// TestWhereClause_NoSelectors tests the case where no selectors are provided.
func TestWhereClause_NoSelectors(t *testing.T) {
	selectors := []clause.Selector{}

	whereClause, whereValues := whereClause(selectors)

	// Expect an empty string and no values since there are no selectors
	assert.Equal(t, "", whereClause)
	assert.Empty(t, whereValues)
}

// TestWhereClause_SingleSelector tests the case where a single selector is
// provided.
func TestWhereClause_SingleSelector(t *testing.T) {
	selectors := []clause.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	whereClause, whereValues := whereClause(selectors)

	expectedClause := "WHERE `user`.`id` = ?"
	assert.Equal(t, expectedClause, whereClause)
	assert.Equal(t, []any{1}, whereValues)
}

// TestWhereClause_MultipleSelectors tests the case where multiple selectors are
// provided.
func TestWhereClause_MultipleSelectors(t *testing.T) {
	selectors := []clause.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
		{Table: "user", Field: "age", Predicate: ">", Value: 18},
	}

	whereClause, whereValues := whereClause(selectors)

	expectedClause := "WHERE `user`.`id` = ? AND `user`.`age` > ?"
	assert.Equal(t, expectedClause, whereClause)
	assert.Equal(t, []any{1, 18}, whereValues)
}

// TestWhereClause_DifferentPredicates tests the case where different predicates
// are provided.
func TestWhereClause_DifferentPredicates(t *testing.T) {
	selectors := []clause.Selector{
		{Table: "user", Field: "name", Predicate: "LIKE", Value: "%Alice%"},
		{Table: "user", Field: "age", Predicate: "<", Value: 30},
	}

	whereClause, whereValues := whereClause(selectors)

	// Expect a WHERE clause with different predicates
	expectedClause := "WHERE `user`.`name` LIKE ? AND `user`.`age` < ?"
	assert.Equal(t, expectedClause, whereClause)
	assert.Equal(t, []any{"%Alice%", 30}, whereValues)
}

// TestBuildBaseGetQuery_NoOptions tests the case where no options are provided.
func TestBuildBaseGetQuery_NoOptions(t *testing.T) {
	getOptions := GetOptions{}

	query, whereValues := Get("user", &getOptions)

	expectedQuery := "SELECT * FROM `user`"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestBuildBaseGetQuery_WithSelectors tests the case where selectors are
// provided.
func TestBuildBaseGetQuery_WithSelectors(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Selectors = []clause.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	query, whereValues := Get("user", &getOptions)

	expectedQuery := "SELECT * FROM `user` WHERE `user`.`id` = ?"
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, []any{1}, whereValues)
}

// TestBuildBaseGetQuery_WithOrders tests the case where orders are provided.
func TestBuildBaseGetQuery_WithOrders(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Orders = []clause.Order{
		{Table: "user", Field: "name", Direction: "ASC"},
	}

	query, whereValues := Get("user", &getOptions)

	expectedQuery := "SELECT * FROM `user` ORDER BY `user`.`name` ASC"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestBuildBaseGetQuery_WithProjections tests the case where projections are
// provided.
func TestBuildBaseGetQuery_WithProjections(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Projections = []clause.Projection{
		{Table: "user", Column: "name", Alias: "user_name"},
	}

	query, whereValues := Get("user", &getOptions)

	expectedQuery := "SELECT `user`.`name` AS `user_name` FROM `user`"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestBuildBaseGetQuery_WithJoins tests the case where joins are provided.
func TestBuildBaseGetQuery_WithJoins(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Joins = []clause.Join{
		{
			Type:  clause.JoinTypeInner,
			Table: "order",
			OnLeft: clause.ColumnSelector{
				Table:   "user",
				Columnn: "id",
			},
			OnRight: clause.ColumnSelector{
				Table:   "order",
				Columnn: "user_id",
			},
		},
	}

	query, whereValues := Get("user", &getOptions)

	expectedQuery := "SELECT * FROM `user` INNER JOIN `order` ON `user`.`id` = `order`.`user_id`"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestBuildBaseGetQuery_WithLock tests the case where the lock option is set.
func TestBuildBaseGetQuery_WithLock(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Lock = true

	query, whereValues := Get("user", &getOptions)

	expectedQuery := "SELECT * FROM `user` FOR UPDATE"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}

// TestBuildBaseGetQuery_WithPage tests the case where pagination is provided.
func TestBuildBaseGetQuery_WithPage(t *testing.T) {
	getOptions := GetOptions{}
	getOptions.Page = &page.Page{Offset: 10, Limit: 20}

	query, whereValues := Get("user", &getOptions)

	expectedQuery := "SELECT * FROM `user` LIMIT 20 OFFSET 10"
	assert.Equal(t, expectedQuery, query)
	assert.Empty(t, whereValues)
}
