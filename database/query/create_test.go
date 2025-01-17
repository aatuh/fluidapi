package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCreateEntity struct {
	ID   int
	Name string
	Age  int
}

// TestGetInsertQueryColumnNames_MultipleColumns tests getInsertQueryColumnNames
// with multiple columns.
func TestGetInsertQueryColumnNames_MultipleColumns(t *testing.T) {
	// Multiple columns
	columns := []string{"id", "name", "age"}

	result := getInsertQueryColumnNames(columns)

	// Expected result
	expectedResult := "`id`, `name`, `age`"

	assert.Equal(t, expectedResult, result)
}

// TestGetInsertQueryColumnNames_SingleColumn tests getInsertQueryColumnNames
// with a single column.
func TestGetInsertQueryColumnNames_SingleColumn(t *testing.T) {
	// Single column
	columns := []string{"id"}

	result := getInsertQueryColumnNames(columns)

	// Expected result
	expectedResult := "`id`"

	assert.Equal(t, expectedResult, result)
}

// TestGetInsertQueryColumnNames_EmptyColumns tests getInsertQueryColumnNames
// with an empty list of columns.
func TestGetInsertQueryColumnNames_EmptyColumns(t *testing.T) {
	// Empty columns
	columns := []string{}

	result := getInsertQueryColumnNames(columns)

	// Expected result is an empty string
	expectedResult := ""

	assert.Equal(t, expectedResult, result)
}

// TestInsertQuery_NormalOperation tests insertQuery with a standard entity.
func TestInsertQuery_NormalOperation(t *testing.T) {
	// Inserter function that returns two columns and values
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		return []string{"id", "name"}, []any{1, "Alice"}
	}

	query, values := Insert(
		&TestCreateEntity{ID: 1, Name: "Alice"},
		"user",
		inserter,
	)

	expectedQuery := "INSERT INTO `user` (`id`, `name`) VALUES (?, ?)"
	expectedValues := []any{1, "Alice"}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsertQuery_SingleColumnEntity tests insertQuery with an entity that has
// only one column.
func TestInsertQuery_SingleColumnEntity(t *testing.T) {
	// Inserter function that returns a single column and value
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		return []string{"id"}, []any{1}
	}

	query, values := Insert(&TestCreateEntity{ID: 1}, "user", inserter)

	expectedQuery := "INSERT INTO `user` (`id`) VALUES (?)"
	expectedValues := []any{1}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}

// TestInsertQuery_NoColumns tests insertQuery with an entity that has no
// columns.
func TestInsertQuery_NoColumns(t *testing.T) {
	// Inserter function that returns no columns or values
	inserter := func(entity *TestCreateEntity) ([]string, []any) {
		return []string{}, []any{}
	}

	query, values := Insert(&TestCreateEntity{}, "user", inserter)

	expectedQuery := "INSERT INTO `user` () VALUES ()"
	expectedValues := []any{}

	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedValues, values)
}
