package test

import (
	"testing"

	"github.com/pakkasys/fluidapi/database"
	"github.com/stretchr/testify/assert"
)

// TestGetByField_Found tests the case where the selector with the given field
// is found.
func TestGetByField_Found(t *testing.T) {
	selectors := database.Selectors{
		{Table: "user", Column: "id", Predicate: "=", Value: 1},
		{Table: "user", Column: "name", Predicate: "=", Value: "Alice"},
	}

	selector := selectors.GetByField("name")

	assert.NotNil(t, selector)
	assert.Equal(t, "user", selector.Table)
	assert.Equal(t, "name", selector.Column)
	assert.Equal(t, "=", string(selector.Predicate))
	assert.Equal(t, "Alice", selector.Value)
}

// TestGetByField_NotFound tests the case where the selector with the given
// field is not found.
func TestGetByField_NotFound(t *testing.T) {
	selectors := database.Selectors{
		{Table: "user", Column: "id", Predicate: "=", Value: 1},
		{Table: "user", Column: "name", Predicate: "=", Value: "Alice"},
	}

	selector := selectors.GetByField("age")

	assert.Nil(t, selector)
}

// TestGetByFields_Found tests the case where the selectors with the given
// fields are found.
func TestGetByFields_Found(t *testing.T) {
	selectors := database.Selectors{
		{Table: "user", Column: "id", Predicate: "=", Value: 1},
		{Table: "user", Column: "name", Predicate: "=", Value: "Alice"},
		{Table: "user", Column: "age", Predicate: ">", Value: 25},
	}

	resultSelectors := selectors.GetByFields("name", "age")

	assert.Len(t, resultSelectors, 2)
	assert.Equal(t, "name", resultSelectors[0].Column)
	assert.Equal(t, "age", resultSelectors[1].Column)
}

// TestGetByFields_NotFound tests the case where none of the selectors with the
// given fields are found.
func TestGetByFields_NotFound(t *testing.T) {
	selectors := database.Selectors{
		{Table: "user", Column: "id", Predicate: "=", Value: 1},
		{Table: "user", Column: "name", Predicate: "=", Value: "Alice"},
	}

	resultSelectors := selectors.GetByFields("age", "address")

	assert.Len(t, resultSelectors, 0)
}

// TestGetByFields_PartialFound tests the case where some selectors with the
// given fields are found.
func TestGetByFields_PartialFound(t *testing.T) {
	selectors := database.Selectors{
		{Table: "user", Column: "id", Predicate: "=", Value: 1},
		{Table: "user", Column: "name", Predicate: "=", Value: "Alice"},
		{Table: "user", Column: "age", Predicate: ">", Value: 25},
	}

	resultSelectors := selectors.GetByFields("name", "address")

	assert.Len(t, resultSelectors, 1)
	assert.Equal(t, "name", resultSelectors[0].Column)
}
